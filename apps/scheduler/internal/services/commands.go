package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/parser"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	commandswithgroupsandresponsesrepository "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	"gorm.io/gorm"
)

type Commands struct {
	db            *gorm.DB
	lock          sync.Mutex
	logger        *slog.Logger
	commandsCache *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
	BusCore       *buscore.Bus
}

func NewCommands(
	db *gorm.DB,
	l *slog.Logger,
	buscore *buscore.Bus,
	commandsRepo commandswithgroupsandresponsesrepository.Repository,
) *Commands {
	return &Commands{
		db:            db,
		logger:        l,
		commandsCache: commandscache.New(commandsRepo, buscore),
		BusCore:       buscore,
	}
}

type channelWithCommandsToCreate struct {
	ChannelID        string         `gorm:"column:channelId" db:"channelId"`
	CommandsToCreate pq.StringArray `gorm:"column:commandsToCreate" db:"commandsToCreate"`
}

func (c *Commands) CreateDefaultCommands(ctx context.Context) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	defaultCommands, err := c.BusCore.Parser.GetDefaultCommands.Request(ctx, struct{}{})
	if err != nil {
		return err
	}

	defaultCommandsNames := make([]string, len(defaultCommands.Data.List))
	for i, command := range defaultCommands.Data.List {
		defaultCommandsNames[i] = command.Name
	}

	var placeholders []string
	var args []interface{}

	for _, name := range defaultCommandsNames {
		placeholders = append(placeholders, "(?)")
		args = append(args, name)
	}

	// "(?), (?), (?)" => "(?), (?), (?)"
	valuesClause := strings.Join(placeholders, ", ")

	query := fmt.Sprintf(
		`
		SELECT
			c.id AS "channelId",
			array_agg(required_commands.name) AS "commandsToCreate"
		FROM
			public.channels c
				CROSS JOIN (
				VALUES %s
			) AS required_commands(name)
				LEFT JOIN
			public.channels_commands cmd
			ON  c.id = cmd."channelId"
				AND required_commands.name = cmd."defaultName"
				AND cmd."default" = true
		WHERE
			cmd.id IS NULL
		GROUP BY
			c.id;
	`, valuesClause,
	)

	var channelsWithCommandsToCreate []channelWithCommandsToCreate
	if err := c.db.
		WithContext(ctx).
		Raw(query, args...).
		Scan(&channelsWithCommandsToCreate).
		Error; err != nil {
		return fmt.Errorf("cannot get channels with commands to create: %w", err)
	}

	if len(channelsWithCommandsToCreate) == 0 {
		return nil
	}

	createModels := make([]model.ChannelsCommands, 0)
	for _, channel := range channelsWithCommandsToCreate {
		var channelRoles []model.ChannelRole
		if err := c.db.Where(
			`"channelId" = ?`,
			channel.ChannelID,
		).Find(&channelRoles).Error; err != nil {
			return fmt.Errorf("cannot get channel roles: %w", err)
		}

		for _, command := range channel.CommandsToCreate {
			defaultCommand, ok := lo.Find(
				defaultCommands.Data.List,
				func(item parser.DefaultCommand) bool {
					return item.Name == command
				},
			)
			if !ok {
				continue
			}

			commandRolesIds := make([]string, 0)
			for _, role := range defaultCommand.RolesNames {
				for _, channelRole := range channelRoles {
					if channelRole.Type == model.ChannelRoleEnum(role) {
						commandRolesIds = append(commandRolesIds, channelRole.ID)
					}
				}
			}

			createModels = append(
				createModels,
				model.ChannelsCommands{
					ID:                        uuid.New().String(),
					Name:                      defaultCommand.Name,
					Cooldown:                  null.IntFrom(0),
					CooldownType:              "GLOBAL",
					Enabled:                   true,
					Aliases:                   defaultCommand.Aliases,
					Description:               null.StringFrom(defaultCommand.Description),
					Visible:                   defaultCommand.Visible,
					ChannelID:                 channel.ChannelID,
					Default:                   true,
					DefaultName:               null.StringFrom(defaultCommand.Name),
					Module:                    defaultCommand.Module,
					IsReply:                   defaultCommand.IsReply,
					KeepResponsesOrder:        defaultCommand.KeepResponsesOrder,
					DeniedUsersIDS:            pq.StringArray{},
					AllowedUsersIDS:           pq.StringArray{},
					RolesIDS:                  commandRolesIds,
					OnlineOnly:                false,
					RequiredWatchTime:         0,
					RequiredMessages:          0,
					RequiredUsedChannelPoints: 0,
				},
			)
		}
	}

	if len(createModels) == 0 {
		return nil
	}

	if err := c.db.WithContext(ctx).CreateInBatches(&createModels, 1000).Error; err != nil {
		return fmt.Errorf("cannot create default commands: %w", err)
	}

	var wg sync.WaitGroup

	for _, channel := range channelsWithCommandsToCreate {
		wg.Add(1)

		go func() {
			defer wg.Done()
			err := c.commandsCache.Invalidate(
				ctx,
				channel.ChannelID,
			)
			if err != nil {
				c.logger.Error(
					"failed to invalidate commands cache",
					logger.Error(err),
					slog.String("channelId", channel.ChannelID),
				)
			}
		}()
	}

	wg.Wait()

	c.logger.Info(
		"Created default commands for channels",
		slog.Int("channels", len(channelsWithCommandsToCreate)),
		slog.Int("commands", len(createModels)),
	)

	return nil
}
