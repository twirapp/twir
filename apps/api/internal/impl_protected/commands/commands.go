package commands

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/commands"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Commands struct {
	*impl_deps.Deps
}

func (c *Commands) convertDbToRpc(cmd *model.ChannelsCommands) *commands.Command {
	return &commands.Command{
		Id:                        cmd.ID,
		Name:                      cmd.Name,
		Cooldown:                  uint32(cmd.Cooldown.Int64),
		CooldownType:              cmd.CooldownType,
		Enabled:                   cmd.Enabled,
		Aliases:                   cmd.Aliases,
		Description:               cmd.Description.String,
		Visible:                   cmd.Visible,
		ChannelId:                 cmd.ChannelID,
		Default:                   cmd.Default,
		DefaultName:               cmd.DefaultName.Ptr(),
		Module:                    cmd.Module,
		IsReply:                   cmd.IsReply,
		KeepResponsesOrder:        cmd.KeepResponsesOrder,
		DeniedUsersIds:            cmd.DeniedUsersIDS,
		AllowedUsersIds:           cmd.AllowedUsersIDS,
		RolesIds:                  cmd.RolesIDS,
		OnlineOnly:                cmd.OnlineOnly,
		RequiredWatchTime:         uint32(cmd.RequiredWatchTime),
		RequiredMessages:          uint32(cmd.RequiredMessages),
		RequiredUsedChannelPoints: uint32(cmd.RequiredUsedChannelPoints),
		Responses: lo.Map(
			cmd.Responses, func(res *model.ChannelsCommandsResponses, _ int) *commands.Command_Response {
				return &commands.Command_Response{
					Id:        res.ID,
					Text:      res.Text.String,
					CommandId: res.CommandID,
					Order:     uint32(res.Order),
				}
			},
		),
		GroupId: cmd.GroupID.Ptr(),
		Group: lo.IfF(
			cmd.Group != nil, func() *commands.Command_Group {
				return &commands.Command_Group{
					Id:        cmd.Group.ID,
					ChannelId: cmd.Group.ChannelID,
					Name:      cmd.Group.Name,
					Color:     cmd.Group.Color,
				}
			},
		).Else(nil),
		CooldownRolesIds: cmd.CooldownRolesIDs,
	}
}

func (c *Commands) CommandsGetAll(
	ctx context.Context,
	_ *emptypb.Empty,
) (*commands.CommandsGetAllResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var cmds []model.ChannelsCommands
	err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ?`, dashboardId).
		Order("name desc").
		Preload("Responses").
		Preload("Group").
		Group(`"id"`).
		Find(&cmds).Error
	if err != nil {
		return nil, err
	}

	return &commands.CommandsGetAllResponse{
		Commands: lo.Map(
			cmds, func(cmd model.ChannelsCommands, _ int) *commands.Command {
				return c.convertDbToRpc(&cmd)
			},
		),
	}, nil
}

func (c *Commands) CommandsGetById(
	ctx context.Context,
	request *commands.GetByIdRequest,
) (*commands.Command, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	cmd := &model.ChannelsCommands{}
	err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "id" = ?`, dashboardId, request.CommandId).
		Preload("Responses").
		Preload("Group").
		Find(cmd).Error
	if err != nil {
		return nil, err
	}
	if cmd.ID == "" {
		return nil, twirp.NewError(twirp.NotFound, "command not found")
	}

	return c.convertDbToRpc(cmd), nil
}

func (c *Commands) CommandsCreate(
	ctx context.Context,
	request *commands.CreateRequest,
) (*commands.Command, error) {
	if len(request.Responses) > 3 {
		return nil, twirp.NewError(twirp.InvalidArgument, "command responses limit is 3")
	}

	dashboardId := ctx.Value("dashboardId").(string)
	command := &model.ChannelsCommands{
		ID:           uuid.New().String(),
		Name:         strings.ToLower(request.Name),
		Cooldown:     null.IntFrom(int64(request.Cooldown)),
		CooldownType: request.CooldownType,
		Enabled:      request.Enabled,
		Aliases: lo.Map(
			lo.IfF(
				request.Aliases == nil, func() []string {
					return []string{}
				},
			).Else(request.Aliases),
			func(alias string, _ int) string {
				return strings.TrimSuffix(strings.ToLower(alias), "!")
			},
		),
		Description:               null.StringFrom(request.Description),
		Visible:                   request.Visible,
		ChannelID:                 dashboardId,
		Default:                   false,
		DefaultName:               null.String{},
		Module:                    "CUSTOM",
		IsReply:                   request.IsReply,
		KeepResponsesOrder:        request.KeepResponsesOrder,
		DeniedUsersIDS:            request.DeniedUsersIds,
		AllowedUsersIDS:           request.AllowedUsersIds,
		RolesIDS:                  request.RolesIds,
		OnlineOnly:                request.OnlineOnly,
		RequiredWatchTime:         int(request.RequiredWatchTime),
		RequiredMessages:          int(request.RequiredMessages),
		RequiredUsedChannelPoints: int(request.RequiredUsedChannelPoints),
		Responses:                 make([]*model.ChannelsCommandsResponses, 0, len(request.Responses)),
		GroupID:                   null.StringFromPtr(request.GroupId),
		CooldownRolesIDs:          request.CooldownRolesIds,
	}

	for _, res := range request.Responses {
		if res.Text == "" {
			continue
		}

		command.Responses = append(
			command.Responses, &model.ChannelsCommandsResponses{
				ID:    uuid.New().String(),
				Text:  null.StringFrom(res.Text),
				Order: int(res.Order),
			},
		)
	}

	err := c.Db.WithContext(ctx).Create(command).Error
	if err != nil {
		return nil, err
	}

	return c.convertDbToRpc(command), nil
}

func (c *Commands) CommandsDelete(
	ctx context.Context,
	request *commands.DeleteRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "id" = ?`, dashboardId, request.CommandId).
		Delete(&model.ChannelsCommands{}).Error
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Commands) CommandsUpdate(
	ctx context.Context,
	request *commands.PutRequest,
) (*commands.Command, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	cmd := &model.ChannelsCommands{}
	err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "id" = ?`, dashboardId, request.Id).
		Find(cmd).Error
	if err != nil {
		return nil, err
	}
	if cmd.ID == "" {
		return nil, twirp.NewError(twirp.NotFound, "command not found")
	}

	if len(request.Command.Responses) > 3 {
		return nil, twirp.NewError(twirp.InvalidArgument, "command responses limit is 3")
	}

	cmd.Name = strings.ToLower(request.Command.Name)
	cmd.Cooldown = null.IntFrom(int64(request.Command.Cooldown))
	cmd.CooldownType = request.Command.CooldownType
	cmd.Enabled = request.Command.Enabled
	cmd.Aliases = lo.Map(
		lo.IfF(
			request.Command.Aliases == nil, func() []string {
				return []string{}
			},
		).Else(request.Command.Aliases),
		func(alias string, _ int) string {
			return strings.TrimSuffix(strings.ToLower(alias), "!")
		},
	)
	cmd.Description = null.StringFrom(request.Command.Description)
	cmd.Visible = request.Command.Visible
	cmd.IsReply = request.Command.IsReply
	cmd.KeepResponsesOrder = request.Command.KeepResponsesOrder
	cmd.AllowedUsersIDS = lo.IfF(
		request.Command.AllowedUsersIds == nil, func() []string {
			return []string{}
		},
	).Else(request.Command.AllowedUsersIds)
	cmd.DeniedUsersIDS = lo.IfF(
		request.Command.DeniedUsersIds == nil, func() []string {
			return []string{}
		},
	).Else(request.Command.DeniedUsersIds)
	cmd.RolesIDS = lo.IfF(
		request.Command.RolesIds == nil, func() []string {
			return []string{}
		},
	).Else(request.Command.RolesIds)
	cmd.OnlineOnly = request.Command.OnlineOnly
	cmd.RequiredWatchTime = int(request.Command.RequiredWatchTime)
	cmd.RequiredMessages = int(request.Command.RequiredMessages)
	cmd.RequiredUsedChannelPoints = int(request.Command.RequiredUsedChannelPoints)
	cmd.GroupID = null.StringFromPtr(request.Command.GroupId)
	cmd.Responses = make([]*model.ChannelsCommandsResponses, 0, len(request.Command.Responses))
	cmd.CooldownRolesIDs = lo.IfF(
		request.Command.CooldownRolesIds == nil, func() []string {
			return []string{}
		},
	).Else(request.Command.CooldownRolesIds)

	for _, res := range request.Command.Responses {
		if res.Text == "" {
			continue
		}

		r := &model.ChannelsCommandsResponses{
			Text:      null.StringFrom(res.Text),
			Order:     int(res.Order),
			CommandID: cmd.ID,
		}

		cmd.Responses = append(cmd.Responses, r)
	}

	txErr := c.Db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			if err = tx.Delete(
				&model.ChannelsCommandsResponses{},
				`"commandId" = ?`,
				cmd.ID,
			).Error; err != nil {
				return err
			}

			return tx.Save(cmd).Error
		},
	)
	if txErr != nil {
		return nil, err
	}

	return c.convertDbToRpc(cmd), nil
}

func (c *Commands) CommandsEnableOrDisable(
	ctx context.Context,
	request *commands.PatchRequest,
) (*commands.Command, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	cmd := &model.ChannelsCommands{}
	err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "id" = ?`, dashboardId, request.CommandId).Find(cmd).Error
	if err != nil {
		return nil, err
	}
	if cmd.ID == "" {
		return nil, twirp.NewError(twirp.NotFound, "command not found")
	}

	cmd.Enabled = request.Enabled
	err = c.Db.WithContext(ctx).Save(cmd).Error
	if err != nil {
		return nil, err
	}

	return c.convertDbToRpc(cmd), nil
}
