package commands_with_groups_and_responses

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/commands"
	commandmodel "github.com/twirapp/twir/libs/repositories/commands/model"
	"github.com/twirapp/twir/libs/repositories/commands_response"
	responsemodel "github.com/twirapp/twir/libs/repositories/commands_response/model"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
)

type UpdateInput struct {
	ActorID   string
	ChannelID string

	Name                      *string
	Cooldown                  *int
	CooldownType              *string
	Enabled                   *bool
	Aliases                   []string
	Description               *string
	Visible                   *bool
	IsReply                   *bool
	KeepResponsesOrder        *bool
	DeniedUsersIDS            []string
	AllowedUsersIDS           []string
	RolesIDS                  []string
	OnlineOnly                *bool
	OfflineOnly               *bool
	CooldownRolesIDs          []string
	EnabledCategories         []string
	RequiredWatchTime         *int
	RequiredMessages          *int
	RequiredUsedChannelPoints *int
	GroupID                   *uuid.UUID
	ExpiresAt                 *time.Time
	ExpiresType               *string
	Responses                 []UpdateInputResponse
}

type UpdateInputResponse struct {
	Text              *string
	Order             int
	TwitchCategoryIDs []string
	OnlineOnly        bool
	OfflineOnly       bool
}

func (c *Service) Update(
	ctx context.Context,
	id uuid.UUID,
	input UpdateInput,
) (entity.CommandWithGroupAndResponses, error) {
	cmds, err := c.commandsWithGroupsAndResponsesRepository.GetManyByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.CommandWithGroupAndResponsesNil, err
	}

	var cmd *model.CommandWithGroupAndResponses
	for _, c := range cmds {
		if c.Command.ID == id {
			cmd = &c
			break
		}
	}

	if cmd == nil {
		return entity.CommandWithGroupAndResponsesNil, fmt.Errorf("command not found")
	}

	if cmd.Command.ChannelID != input.ChannelID {
		return entity.CommandWithGroupAndResponsesNil, fmt.Errorf("command not found")
	}

	if cmd.Command.Default && input.ExpiresType != nil && *input.ExpiresType == "DELETE" {
		return entity.CommandWithGroupAndResponsesNil, fmt.Errorf("default command cannot be deleted")
	}

	onlyCmds := make([]commandmodel.Command, 0, len(cmds))
	for _, c := range cmds {
		onlyCmds = append(onlyCmds, c.Command)
	}

	if input.Name != nil {
		if conflict, _ := c.commandsService.IsNameConflicting(
			onlyCmds,
			*input.Name,
			input.Aliases,
			[]uuid.UUID{cmd.Command.ID},
		); conflict {
			return entity.CommandWithGroupAndResponsesNil, fmt.Errorf("command with this name or alias already exists")
		}
	}

	if input.Aliases != nil {
		if conflict, _ := c.commandsService.IsNameConflicting(
			onlyCmds,
			cmd.Command.Name,
			input.Aliases,
			[]uuid.UUID{cmd.Command.ID},
		); conflict {
			return entity.CommandWithGroupAndResponsesNil, fmt.Errorf("command with this name or alias already exists")
		}
	}

	commandUpdateInput := commands.UpdateInput{
		Name:                      input.Name,
		Cooldown:                  input.Cooldown,
		CooldownType:              input.CooldownType,
		Enabled:                   input.Enabled,
		Aliases:                   input.Aliases,
		Description:               input.Description,
		Visible:                   input.Visible,
		IsReply:                   input.IsReply,
		KeepResponsesOrder:        input.KeepResponsesOrder,
		DeniedUsersIDS:            input.DeniedUsersIDS,
		AllowedUsersIDS:           input.AllowedUsersIDS,
		RolesIDS:                  input.RolesIDS,
		OnlineOnly:                input.OnlineOnly,
		OfflineOnly:               input.OfflineOnly,
		CooldownRolesIDs:          input.CooldownRolesIDs,
		EnabledCategories:         input.EnabledCategories,
		RequiredWatchTime:         input.RequiredWatchTime,
		RequiredMessages:          input.RequiredMessages,
		RequiredUsedChannelPoints: input.RequiredUsedChannelPoints,
		GroupID:                   input.GroupID,
		ExpiresAt:                 input.ExpiresAt,
		ExpiresType:               input.ExpiresType,
	}

	var newCmd model.CommandWithGroupAndResponses
	trErr := c.trmManager.Do(
		ctx,
		func(trCtx context.Context) error {
			newDbCmd, err := c.commandsRepository.Update(trCtx, id, commandUpdateInput)
			if err != nil {
				return err
			}

			newCmd.Command = newDbCmd

			if input.Responses != nil {
				for _, r := range cmd.Responses {
					err := c.responsesRepository.Delete(trCtx, r.ID)
					if err != nil {
						return err
					}
				}

				newCmd.Responses = make([]responsemodel.Response, 0, len(cmd.Responses))
				for _, r := range input.Responses {
					newResponse, err := c.responsesRepository.Create(
						trCtx,
						commands_response.CreateInput{
							CommandID:         newDbCmd.ID,
							Text:              r.Text,
							Order:             r.Order,
							TwitchCategoryIDs: r.TwitchCategoryIDs,
							OnlineOnly:        r.OnlineOnly,
							OfflineOnly:       r.OfflineOnly,
						},
					)
					if err != nil {
						return err
					}

					newCmd.Responses = append(newCmd.Responses, newResponse)
				}
			}

			return nil
		},
	)
	if trErr != nil {
		return entity.CommandWithGroupAndResponsesNil, trErr
	}

	if err := c.cachedCommandsClient.Invalidate(ctx, input.ChannelID); err != nil {
		c.logger.Error("failed to invalidate cached commands", err)
	}

	c.logger.Audit(
		"Command edited",
		audit.Fields{
			OldValue:      cmd,
			NewValue:      newCmd,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelCommand),
			OperationType: audit.OperationUpdate,
			ObjectID:      lo.ToPtr(newCmd.Command.ID.String()),
		},
	)

	return c.mapToEntity(newCmd), nil
}
