package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_responses"
	"github.com/twirapp/twir/libs/audit"
	commandwithrelationentity "github.com/twirapp/twir/libs/entities/command_with_relations"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/errors"
	"github.com/twirapp/twir/libs/repositories/command_role_cooldown"
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
	EnabledCategories         []string
	RequiredWatchTime         int
	RequiredMessages          int
	RequiredUsedChannelPoints int
	GroupID                   *uuid.UUID
	ExpiresAt                 *int
	ExpiresType               *string
	Responses                 []CreateInputResponse
	RoleCooldowns             []CreateInputRoleCooldown
	Platforms                 []platform.Platform
}

type CreateInputResponse struct {
	Text              *string
	Order             int
	TwitchCategoryIDs []string
	OnlineOnly        bool
	OfflineOnly       bool
}

type CreateInputRoleCooldown struct {
	RoleID   string
	Cooldown int
}

func (c *Service) Create(ctx context.Context, input CreateInput) (commandwithrelationentity.Command, error) {
	plan, err := c.plansRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return commandwithrelationentity.CommandNil, errors.NewInternalError("failed to get plan", err)
	}
	if plan.IsNil() {
		return commandwithrelationentity.CommandNil, errors.NewNotFoundError("Plan configuration not found for your channel")
	}

	cmds, err := c.commandsRepository.GetManyByChannelID(ctx, input.ChannelID)
	if err != nil {
		return commandwithrelationentity.CommandNil, errors.NewInternalError("failed to get commands", err)
	}

	var createdCommands int
	for _, cmd := range cmds {
		if !cmd.Default {
			createdCommands++
		}
	}

	if createdCommands >= plan.MaxCommands {
		return commandwithrelationentity.CommandNil, errors.NewBadRequestError(
			fmt.Sprintf("You have reached the maximum limit of %d commands", plan.MaxCommands),
		)
	}

	if len(input.Responses) > plan.MaxCommandsResponses {
		return commandwithrelationentity.CommandNil, errors.NewBadRequestError(
			fmt.Sprintf("You have reached the maximum limit of %d responses per command", plan.MaxCommandsResponses),
		)
	}

	isNameConflict, err := c.IsNameConflicting(
		cmds,
		input.Name,
		input.Aliases,
		nil,
	)
	if err != nil {
		return commandwithrelationentity.CommandNil, errors.NewInternalError("failed to check name conflict", err)
	}
	if isNameConflict {
		return commandwithrelationentity.CommandNil, errors.NewConflictError("A command with this name or alias already exists")
	}

	aliases := make([]string, 0, len(input.Aliases))
	for _, alias := range input.Aliases {
		alias = strings.TrimSuffix(strings.ToLower(alias), "!")
		if alias != "" {
			aliases = append(aliases, alias)
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
					EnabledCategories:         input.EnabledCategories,
					RequiredWatchTime:         input.RequiredWatchTime,
					RequiredMessages:          input.RequiredMessages,
					RequiredUsedChannelPoints: input.RequiredUsedChannelPoints,
					GroupID:                   input.GroupID,
					ExpiresAt:                 expiresAt,
					ExpiresType:               input.ExpiresType,
					Platforms:                 input.Platforms,
				},
			)
			if err != nil {
				return errors.NewInternalError("failed to create command", err)
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
					return errors.NewInternalError("failed to create command response", err)
				}
			}

			for _, roleCooldown := range input.RoleCooldowns {
				roleID, err := uuid.Parse(roleCooldown.RoleID)
				if err != nil {
					return errors.NewBadRequestError(fmt.Sprintf("invalid role ID: %s", roleCooldown.RoleID))
				}

				_, err = c.commandRoleCooldownsRepository.Create(
					txCtx,
					command_role_cooldown.CreateInput{
						CommandID: dbCmd.ID,
						RoleID:    roleID,
						Cooldown:  roleCooldown.Cooldown,
					},
				)
				if err != nil {
					return errors.NewInternalError("failed to create role cooldown", err)
				}
			}

			if err := c.cachedCommandsClient.Invalidate(ctx, input.ChannelID); err != nil {
				return errors.NewInternalError("failed to invalidate cached commands", err)
			}

			return nil
		},
	)
	if trErr != nil {
		return commandwithrelationentity.CommandNil, trErr
	}

	convertedCommand := c.modelToEntity(dbCmd)

	_ = c.auditRecorder.RecordCreateOperation(
		ctx,
		audit.CreateOperation{
			Metadata: audit.OperationMetadata{
				System:    "channels_commands",
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(convertedCommand.ID.String()),
			},
			NewValue: convertedCommand,
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
