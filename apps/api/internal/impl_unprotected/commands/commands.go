package commands

import (
	"context"
	"go.uber.org/zap"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/commands_unprotected"
)

type Commands struct {
	*impl_deps.Deps
}

func (c *Commands) GetChannelCommands(
	ctx context.Context,
	req *commands_unprotected.GetChannelCommandsRequest,
) (*commands_unprotected.GetChannelCommandsResponse, error) {
	var entities []*model.ChannelsCommands
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND visible = ? AND enabled = ?`, req.ChannelId, true, true).
		Preload("Responses").
		Preload("Group").
		Find(&entities).Error; err != nil {
		return nil, err
	}

	return &commands_unprotected.GetChannelCommandsResponse{
		Commands: lo.Map(
			entities,
			func(cmd *model.ChannelsCommands, _ int) *commands_unprotected.Command {
				var roles []*model.ChannelRole
				if len(cmd.RolesIDS) > 0 {
					ids := lo.Map(cmd.RolesIDS, func(id string, _ int) string { return id })
					err := c.Db.WithContext(ctx).
						Where(`"channelId" = ? AND "id" IN ?`, req.ChannelId, ids).
						Find(&roles).Error

					if err != nil {
						zap.S().Error(err)
					}
				}

				return &commands_unprotected.Command{
					Name: cmd.Name,
					Responses: lo.Map(
						cmd.Responses,
						func(item *model.ChannelsCommandsResponses, _ int) string {
							return item.Text.String
						},
					),
					Cooldown:     cmd.Cooldown.Int64,
					CooldownType: cmd.CooldownType,
					Aliases:      cmd.Aliases,
					Description:  cmd.Description.Ptr(),
					Group: lo.IfF(
						cmd.Group != nil, func() *string {
							return &cmd.Group.Name
						},
					).Else(nil),
					Module: &cmd.Module,
					Permissions: lo.Map(
						roles,
						func(role *model.ChannelRole, _ int) *commands_unprotected.Command_Permission {
							return &commands_unprotected.Command_Permission{
								Name: role.Name,
								Type: role.Type.String(),
							}
						},
					),
				}
			},
		),
	}, nil
}
