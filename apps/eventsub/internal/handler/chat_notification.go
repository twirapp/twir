package handler

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/bus-core/events"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
	"github.com/twirapp/twitchy/eventsub"
)

func (c *Handler) HandleChannelChatNotification(
	ctx context.Context,
	event eventsub.ChannelChatNotificationEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	switch event.NoticeType {
	case "sub_gift":
		c._notificationSubGift(ctx, event)
	case "community_sub_gift":
		c._notificationCommunitySubGift(ctx, event)
	}
}

func (c *Handler) _notificationSubGift(
	ctx context.Context,
	event eventsub.ChannelChatNotificationEvent,
) {
	if event.SubGift == nil {
		return
	}

	tier := getSubPlan(string(event.SubGift.SubTier))

	c.logger.Info(
		"subgift",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("targetUserName", event.SubGift.RecipientUserName),
		slog.String("targetUserId", event.SubGift.RecipientUserId),
		slog.String("level", string(event.SubGift.SubTier)),
	)

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserId,
			Type:      model.ChannelEventListItemTypeSubGift,
			Data: &model.ChannelsEventsListItemData{
				SubGiftUserName:              event.ChatterUserLogin,
				SubGiftUserDisplayName:       event.ChatterUserName,
				SubGiftTargetUserName:        event.SubGift.RecipientUserName,
				SubGiftTargetUserDisplayName: event.SubGift.RecipientUserName,
			},
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	c.twirBus.Events.SubGift.Publish(
		ctx,
		events.SubGiftMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			SenderUserID:      event.ChatterUserId,
			SenderUserName:    event.ChatterUserLogin,
			SenderDisplayName: event.ChatterUserName,
			TargetUserName:    event.SubGift.RecipientUserName,
			TargetDisplayName: event.SubGift.RecipientUserName,
			Level:             tier,
		},
	)
}

func (c *Handler) _notificationCommunitySubGift(
	ctx context.Context,
	event eventsub.ChannelChatNotificationEvent,
) {
}
