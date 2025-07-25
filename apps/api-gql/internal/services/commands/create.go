package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_responses"
	"github.com/twirapp/twir/libs/repositories/commands"
	"github.com/twirapp/twir/libs/repositories/commands/model"
)

type CreateInput struct {
	ChannelID string
	ActorID   string

	Name                      string
	Cooldown                  int
	CooldownType              string
	Enabled                   bool
	Aliases                   []string
	Description               string
	Visible                   bool
	IsReply                   bool
	KeepResponsesOrder        bool
	DeniedUsersIDS            []string
	AllowedUsersIDS           []string
	RolesIDS                  []string
	OnlineOnly                bool
	CooldownRolesIDs          []string
	EnabledCategories         []string
	RequiredWatchTime         int
	RequiredMessages          int
	RequiredUsedChannelPoints int
	GroupID                   *uuid.UUID
	ExpiresAt                 *int
	ExpiresType               *string
	Responses                 []CreateInputResponse
}

type CreateInputResponse struct {
	Text              *string
	Order             int
	TwitchCategoryIDs []string
	OnlineOnly        bool
	OfflineOnly       bool
}

func (c *Service) Create(ctx context.Context, input CreateInput) (entity.Command, error) {
	cmds, err := c.commandsRepository.GetManyByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.CommandNil, fmt.Errorf("failed to get commands: %w", err)
	}

	var createdCommands int
	for _, cmd := range cmds {
		if !cmd.Default {
			createdCommands++
		}
	}

	if createdCommands >= maxCommands {
		return entity.CommandNil, fmt.Errorf("maximum commands limit reached")
	}

	isNameConflict, err := c.IsNameConflicting(
		cmds,
		input.Name,
		input.Aliases,
		nil,
	)
	if err != nil {
		return entity.CommandNil, err
	}
	if isNameConflict {
		return entity.CommandNil, fmt.Errorf("command with this name or aliase already exists")
	}

	aliases := make([]string, 0, len(input.Aliases))
	for _, aliase := range input.Aliases {
		aliase = strings.TrimSuffix(strings.ToLower(aliase), "!")
		if aliase != "" {
			aliases = append(aliases, aliase)
		}
	}

	var expiresAt *time.Time
	if input.ExpiresAt != nil {
		expiresAt = lo.ToPtr(time.UnixMilli(int64(*input.ExpiresAt)))
	}

	var dbCmd model.Command
	trErr := c.trManager.Do(
		ctx,
		func(txCtx context.Context) error {
			newCommand, err := c.commandsRepository.Create(
				txCtx,
				commands.CreateInput{
					ChannelID:                 input.ChannelID,
					Name:                      input.Name,
					Cooldown:                  input.Cooldown,
					CooldownType:              input.CooldownType,
					Enabled:                   input.Enabled,
					Aliases:                   aliases,
					Description:               input.Description,
					Visible:                   input.Visible,
					IsReply:                   input.IsReply,
					KeepResponsesOrder:        input.KeepResponsesOrder,
					DeniedUsersIDS:            input.DeniedUsersIDS,
					AllowedUsersIDS:           input.AllowedUsersIDS,
					RolesIDS:                  input.RolesIDS,
					OnlineOnly:                input.OnlineOnly,
					CooldownRolesIDs:          input.CooldownRolesIDs,
					EnabledCategories:         input.EnabledCategories,
					RequiredWatchTime:         input.RequiredWatchTime,
					RequiredMessages:          input.RequiredMessages,
					RequiredUsedChannelPoints: input.RequiredUsedChannelPoints,
					GroupID:                   input.GroupID,
					ExpiresAt:                 expiresAt,
					ExpiresType:               input.ExpiresType,
				},
			)
			if err != nil {
				return fmt.Errorf("failed to create command: %w", err)
			}

			dbCmd = newCommand

			for _, response := range input.Responses {
				_, err := c.commandsResponsesService.Create(
					txCtx,
					commands_responses.CreateInput{
						CommandID:         dbCmd.ID,
						Text:              response.Text,
						Order:             response.Order,
						TwitchCategoryIDs: response.TwitchCategoryIDs,
						OnlineOnly:        response.OnlineOnly,
						OfflineOnly:       response.OfflineOnly,
					},
				)
				if err != nil {
					return fmt.Errorf("failed to create command response: %w", err)
				}
			}

			if err := c.cachedCommandsClient.Invalidate(ctx, input.ChannelID); err != nil {
				return fmt.Errorf("failed to invalidate cached commands: %w", err)
			}

			return nil
		},
	)
	if trErr != nil {
		return entity.CommandNil, fmt.Errorf("failed to create command: %w", trErr)
	}

	convertedCommand := c.modelToEntity(dbCmd)

	c.logger.Audit(
		"New command created",
		audit.Fields{
			NewValue:      convertedCommand,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			System:        "channels_commands", // TODO: use some enum
			OperationType: audit.OperationCreate,
			ObjectID:      lo.ToPtr(convertedCommand.ID.String()),
		},
	)

	return convertedCommand, nil
}

func (c *Service) CreateMultiple(ctx context.Context, input []CreateInput) error {
	txErr := c.trManager.Do(
		ctx,
		func(txCtx context.Context) error {
			for _, cmd := range input {
				_, err := c.Create(txCtx, cmd)
				if err != nil {
					return err
				}
			}

			return nil
		},
	)
	if txErr != nil {
		return fmt.Errorf("failed to create commands: %w", txErr)
	}

	return nil
}
