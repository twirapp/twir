package grpc

import (
	"context"
	"github.com/satont/tsuwari/apps/events/internal"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/events"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EventsGrpcImplementation struct {
	events.UnimplementedEventsServer

	services *internal.Services
}

func NewEvents(services *internal.Services) *EventsGrpcImplementation {
	return &EventsGrpcImplementation{
		services: services,
	}
}

func (c *EventsGrpcImplementation) Follow(_ context.Context, msg *events.FollowMessage) (*emptypb.Empty, error) {
	dbEntity := model.Event{}

	c.processOperations(dbEntity.Operations, Data{
		UserName:    msg.UserName,
		RaidViewers: 0,
	})

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) Subscribe(context.Context, *events.SubscribeMessage) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) ReSubscribe(context.Context, *events.ReSubscribeMessage) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) RedemptionCreated(context.Context, *events.RedemptionCreatedMessage) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) CommandUsed(context.Context, *events.CommandUsedMessage) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) FirstUserMessage(context.Context, *events.FirstUserMessageMessage) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) Raided(context.Context, *events.RaidedMessage) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) TitleOrCategoryChanged(context.Context, *events.TitleOrCategoryChangedMessage) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) StreamOnline(context.Context, *events.StreamOnlineMessage) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) StreamOffline(context.Context, *events.StreamOfflineMessage) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
