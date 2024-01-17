package grpm_impl

import (
	"context"

	"github.com/satont/twir/apps/eventsub/internal/client"
	"github.com/satont/twir/apps/eventsub/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/eventsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EventSubGrpcImpl struct {
	eventsub.UnimplementedEventSubServer
	eventSubClient *client.SubClient
	callbackUrl    string
	services       *types.Services
}

func NewGrpcImpl(
	eventSubClient *client.SubClient,
	services *types.Services,
	callBackUrl string,
) *EventSubGrpcImpl {
	return &EventSubGrpcImpl{
		eventSubClient: eventSubClient,
		callbackUrl:    callBackUrl,
		services:       services,
	}
}

func (c *EventSubGrpcImpl) SubscribeToEvents(
	ctx context.Context, msg *eventsub.SubscribeToEventsRequest,
) (*emptypb.Empty, error) {
	channel := model.Channels{}
	err := c.services.Gorm.Where(
		`"id" = ?`,
		msg.ChannelId,
	).First(&channel).Error
	if err != nil {
		return nil, status.Error(codes.NotFound, "channel not found")
	}

	if err := c.eventSubClient.SubscribeToNeededEvents(
		ctx,
		msg.ChannelId,
		channel.BotID,
	); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
