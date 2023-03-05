package grpm_impl

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/eventsub"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EventSubGrpcImpl struct {
	eventsub.UnimplementedEventSubServer
}

func NewGrpcImpl() *EventSubGrpcImpl {
	return &EventSubGrpcImpl{}
}

func (c *EventSubGrpcImpl) SubscribeToEvents(_ context.Context, msg *eventsub.SubscribeToEventsRequest) (*emptypb.Empty, error) {
	panic("implement me")
}
