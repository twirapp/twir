package grpc_impl

import (
	"context"
	"github.com/satont/tsuwari/apps/events/internal"
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
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
		},
		"FOLLOW",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) Subscribe(_ context.Context, msg *events.SubscribeMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			SubLevel:        msg.Level,
		},
		"SUBSCRIBE",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) ReSubscribe(_ context.Context, msg *events.ReSubscribeMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			SubLevel:        msg.Level,
			ResubMessage:    msg.Message,
			ResubMonths:     msg.Months,
			ResubStreak:     msg.Streak,
		},
		"RESUBSCRIBE",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) RedemptionCreated(_ context.Context, msg *events.RedemptionCreatedMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			RewardCost:      msg.RewardCost,
			RewardInput:     msg.Input,
			RewardName:      msg.RewardName,
		},
		"REDEMPTION_CREATED",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) CommandUsed(_ context.Context, msg *events.CommandUsedMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			CommandName:     msg.CommandName,
		},
		"COMMAND_USED",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) FirstUserMessage(_ context.Context, msg *events.FirstUserMessageMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
		},
		"FIRST_USER_MESSAGE",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) Raided(_ context.Context, msg *events.RaidedMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			RaidViewers:     msg.Viewers,
		},
		"RAIDED",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) TitleOrCategoryChanged(_ context.Context, msg *events.TitleOrCategoryChangedMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		Data{
			OldCategory: msg.OldCategory,
			NewCategory: msg.NewCategory,
			OldTitle:    msg.OldTitle,
			NewTitle:    msg.NewTitle,
		},
		"TITLE_OR_CATEGORY_CHANGED",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) StreamOnline(_ context.Context, msg *events.StreamOnlineMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		Data{
			StreamTitle:    msg.Title,
			StreamCategory: msg.Category,
		},
		"STREAM_ONLINE",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) StreamOffline(_ context.Context, msg *events.StreamOfflineMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		Data{},
		"STREAM_OFFLINE",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) SubGift(_ context.Context, msg *events.SubGiftMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		Data{
			TargetUserName:        msg.TargetUserName,
			TargetUserDisplayName: msg.TargetDisplayName,
			SubLevel:              msg.Level,
		},
		"SUB_GIFT",
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) ChatClear(_ context.Context, msg *events.ChatClearMessage) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		Data{},
		"CHAT_CLEAR",
	)

	return &emptypb.Empty{}, nil
}
