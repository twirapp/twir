package commands

import (
	"context"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/api/commands"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
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
		RequiredWatchTime:         uint64(cmd.RequiredWatchTime),
		RequiredMessages:          uint64(cmd.RequiredMessages),
		RequiredUsedChannelPoints: uint64(cmd.RequiredUsedChannelPoints),
		Responses: lo.Map(cmd.Responses, func(res *model.ChannelsCommandsResponses, _ int) *commands.Command_Response {
			return &commands.Command_Response{
				Id:        res.ID,
				Text:      res.Text.Ptr(),
				CommandId: res.CommandID,
				Order:     uint32(res.Order),
			}
		}),
		GroupId: cmd.GroupID.Ptr(),
	}
}

func (c *Commands) CommandsGetAll(ctx context.Context, empty *emptypb.Empty) (*commands.CommandsGetAllResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var cmds []model.ChannelsCommands
	err := c.Db.
		Where(`"channelId" = ?`, dashboardId).
		Preload("Responses").
		Preload("Group").
		Find(&cmds).Error
	if err != nil {
		return nil, err
	}

	return &commands.CommandsGetAllResponse{
		Commands: lo.Map(cmds, func(cmd model.ChannelsCommands, _ int) *commands.Command {
			return c.convertDbToRpc(&cmd)
		}),
	}, nil
}

func (c *Commands) CommandsGetById(ctx context.Context, request *commands.GetByIdRequest) (*commands.Command, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	cmd := &model.ChannelsCommands{}
	err := c.Db.
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

func (c *Commands) CommandsCreate(ctx context.Context, request *commands.CreateRequest) (*commands.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Commands) CommandsDelete(ctx context.Context, request *commands.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Commands) CommandsUpdate(ctx context.Context, request *commands.PutRequest) (*commands.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Commands) CommandsEnableOrDisable(ctx context.Context, request *commands.PatchRequest) (*commands.Command, error) {
	//TODO implement me
	panic("implement me")
}
