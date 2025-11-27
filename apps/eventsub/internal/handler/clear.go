package handler

import (
	"context"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/logger"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
)

func (c *Handler) HandleChannelChatClear(
	ctx context.Context,
	event eventsub.ChannelChatClearEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel chat clear",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("channelName", event.BroadcasterUserLogin),
	)

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserId,
			UserID:    nil,
			Type:      model.ChannelEventListItemTypeChatClear,
			Data:      &model.ChannelsEventsListItemData{},
		},
	); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}

	c.twirBus.Events.ChatClear.Publish(
		ctx,
		events.ChatClearMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
		},
	)
}
