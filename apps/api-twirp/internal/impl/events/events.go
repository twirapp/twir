package events

import (
	"context"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/deps"
	"github.com/satont/tsuwari/libs/grpc/generated/api/events"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Events struct {
	*deps.Deps
}

func (c *Events) EventsGetAll(ctx context.Context, empty *emptypb.Empty) (*events.GetAllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Events) EventsGetById(ctx context.Context, request *events.GetByIdRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Events) EventsCreate(ctx context.Context, request *events.CreateRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Events) EventsDelete(ctx context.Context, request *events.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Events) EventsPut(ctx context.Context, request *events.PutRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Events) EventsPatch(ctx context.Context, request *events.PatchRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}
