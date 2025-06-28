package handler

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/bus-core/events"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelChatNotification(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelChatNotification,
) {
	switch event.NoticeType {
	case "sub_gift":
		c._notificationSubGift(event)
	case "community_sub_gift":
		c._notificationCommunitySubGift(event)
	}
}

func (c *Handler) _notificationSubGift(event *eventsub_bindings.EventChannelChatNotification) {
	if event.SubGift == nil {
		return
	}

	tier := getSubPlan(event.SubGift.SubTier)

	c.logger.Info(
		"subgift",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("targetUserName", event.SubGift.RecipientUserName),
		slog.String("targetUserId", event.SubGift.RecipientUserID),
		slog.String("level", event.SubGift.SubTier),
	)

	if err := c.eventsListRepository.Create(
		context.TODO(),
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserID,
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
		events.SubGiftMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
				ChannelName: event.BroadcasterUserLogin,
			},
			SenderUserID:      event.ChatterUserID,
			SenderUserName:    event.ChatterUserLogin,
			SenderDisplayName: event.ChatterUserName,
			TargetUserName:    event.SubGift.RecipientUserName,
			TargetDisplayName: event.SubGift.RecipientUserName,
			Level:             tier,
		},
	)
}

func (c *Handler) _notificationCommunitySubGift(
	event *eventsub_bindings.
		EventChannelChatNotification,
) {

}
