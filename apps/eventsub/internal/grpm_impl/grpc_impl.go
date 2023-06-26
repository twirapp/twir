package grpm_impl

import (
	"context"
	"github.com/satont/twir/apps/eventsub/internal/client"
	"github.com/satont/twir/apps/eventsub/internal/types"
	"github.com/satont/twir/libs/grpc/generated/eventsub"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EventSubGrpcImpl struct {
	eventsub.UnimplementedEventSubServer
	eventSubClient *client.SubClient
	callbackUrl    string
	services       *types.Services
}

func NewGrpcImpl(eventSubClient *client.SubClient, services *types.Services, callBackUrl string) *EventSubGrpcImpl {
	return &EventSubGrpcImpl{
		eventSubClient: eventSubClient,
		callbackUrl:    callBackUrl,
		services:       services,
	}
}

func (c *EventSubGrpcImpl) SubscribeToEvents(
	ctx context.Context, msg *eventsub.SubscribeToEventsRequest,
) (*emptypb.Empty, error) {
	err := c.eventSubClient.SubscribeToNeededEvents(ctx, msg.ChannelId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
