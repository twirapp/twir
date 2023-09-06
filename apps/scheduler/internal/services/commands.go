package services

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Commands struct {
	db         *gorm.DB
	parserGrpc parser.ParserClient
	lock       sync.Mutex
}

func NewCommands(db *gorm.DB, parserGrpc parser.ParserClient) *Commands {
	return &Commands{
		db:         db,
		parserGrpc: parserGrpc,
	}
}

func (c *Commands) CreateDefaultCommands(ctx context.Context, usersIds []string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	defaultCommands, err := c.parserGrpc.GetDefaultCommands(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	var channelsWithCommands []model.Channels
	if err := c.db.
		WithContext(ctx).
		Preload("Commands", c.db.Where(&model.ChannelsCommands{Default: true})).
		Find(&channelsWithCommands).
		Where(`"id" IN ?`, usersIds).
		Error; err != nil {
		return err
	}

	for _, channel := range channelsWithCommands {
		for _, command := range defaultCommands.List {
			// skip if command exists
			if lo.SomeBy(
				channel.Commands,
				func(c model.ChannelsCommands) bool {
					return c.DefaultName.String == command.Name
				},
			) {
				continue
			}

			var channelRoles []model.ChannelRole
			if err := c.db.Where(`"channelId" = ?`, channel.ID).Find(&channelRoles).Error; err != nil {
				return err
			}

			commandRolesIds := make([]string, 0, len(command.RolesNames))
			for _, role := range command.RolesNames {
				for _, channelRole := range channelRoles {
					if channelRole.Type == model.ChannelRoleEnum(role) {
						commandRolesIds = append(commandRolesIds, channelRole.ID)
					}
				}
			}

			newCmd := &model.ChannelsCommands{
				ID:                        uuid.New().String(),
				Name:                      command.Name,
				Cooldown:                  null.IntFrom(0),
				CooldownType:              "GLOBAL",
				Enabled:                   true,
				Aliases:                   command.Aliases,
				Description:               null.StringFrom(command.Description),
				Visible:                   command.Visible,
				ChannelID:                 channel.ID,
				Default:                   true,
				DefaultName:               null.StringFrom(command.Name),
				Module:                    command.Module,
				IsReply:                   command.IsReply,
				KeepResponsesOrder:        command.KeepResponsesOrder,
				DeniedUsersIDS:            pq.StringArray{},
				AllowedUsersIDS:           pq.StringArray{},
				RolesIDS:                  commandRolesIds,
				OnlineOnly:                false,
				RequiredWatchTime:         0,
				RequiredMessages:          0,
				RequiredUsedChannelPoints: 0,
			}

			if err := c.db.
				WithContext(ctx).
				Create(newCmd).
				Error; err != nil {
				zap.S().Error("cannot create default command", zap.Error(err), zap.Any("command", newCmd))
				return err
			}
		}
	}

	return nil
}
