package commands

import (
	"context"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/deps"
	"github.com/satont/tsuwari/libs/grpc/generated/api/commands"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Commands struct {
	*deps.Deps
}

func (c *Commands) CommandsGetAll(ctx context.Context, empty *emptypb.Empty) (*commands.CommandsGetAllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Commands) CommandsGetById(ctx context.Context, request *commands.GetByIdRequest) (*commands.Command, error) {
	//TODO implement me
	panic("implement me")
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
