package handler

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/bus-core/events"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelChatClear(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelChatClear,
) {
	c.logger.Info(
		"channel chat clear",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
	)

	if err := c.eventsListRepository.Create(
		context.TODO(),
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserID,
			UserID:    nil,
			Type:      model.ChannelEventListItemTypeChatClear,
			Data:      &model.ChannelsEventsListItemData{},
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	c.twirBus.Events.ChatClear.Publish(
		events.ChatClearMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
				ChannelName: event.BroadcasterUserLogin,
			},
		},
	)
}
