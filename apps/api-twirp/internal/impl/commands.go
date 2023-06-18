package impl

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/api/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *Api) CommandsCreate(ctx context.Context, request *api.CreateRequest) (*api.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) CommandsGetById(ctx context.Context, request *api.GetByIdRequest) (*api.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) CommandsGetAll(ctx context.Context, request *api.GetAllRequest) (*api.GetAllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) CommandsDelete(ctx context.Context, request *api.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) CommandsPut(ctx context.Context, request *api.CreateRequest) (*api.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) CommandsPatch(ctx context.Context, request *api.PatchRequest) (*api.Command, error) {
	//TODO implement me
	panic("implement me")
}
