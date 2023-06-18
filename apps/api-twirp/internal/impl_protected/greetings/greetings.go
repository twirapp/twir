package greetings

import (
	"context"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	"github.com/satont/tsuwari/libs/grpc/generated/api/greetings"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Greetings struct {
	*impl_deps.Deps
}

func (c *Greetings) GreetingsGetAll(ctx context.Context, empty *emptypb.Empty) (*greetings.GetAllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Greetings) GreetingsGetById(ctx context.Context, request *greetings.GetByIdRequest) (*greetings.Greeting, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Greetings) GreetingsCreate(ctx context.Context, request *greetings.CreateRequest) (*greetings.Greeting, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Greetings) GreetingsDelete(ctx context.Context, request *greetings.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Greetings) GreetingsUpdate(ctx context.Context, request *greetings.PutRequest) (*greetings.Greeting, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Greetings) GreetingsEnableOrDisable(ctx context.Context, request *greetings.PatchRequest) (*greetings.Greeting, error) {
	//TODO implement me
	panic("implement me")
}
