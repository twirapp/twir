package handler

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/twirapp/twir/libs/bus-core/events"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelRaid(
	ctx context.Context,
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelRaid,
) {
	c.logger.Info(
		"channel raid",
		slog.String("channelId", event.ToBroadcasterUserID),
		slog.String("channelName", event.ToBroadcasterUserName),
		slog.String("userId", event.FromBroadcasterUserID),
		slog.String("userName", event.FromBroadcasterUserLogin),
		slog.Int("viewers", event.Viewers),
	)

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.ToBroadcasterUserID,
			UserID:    &event.FromBroadcasterUserID,
			Type:      model.ChannelEventListItemTypeRaided,
			Data: &model.ChannelsEventsListItemData{
				RaidedViewersCount:    strconv.Itoa(event.Viewers),
				RaidedFromUserName:    event.FromBroadcasterUserLogin,
				RaidedFromDisplayName: event.FromBroadcasterUserName,
			},
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	if err := c.twirBus.Events.Raided.Publish(
		ctx,
		events.RaidedMessage{
			BaseInfo: events.BaseInfo{
				ChannelName: event.ToBroadcasterUserID,
				ChannelID:   event.ToBroadcasterUserID,
			},
			UserID:          event.FromBroadcasterUserID,
			UserName:        event.FromBroadcasterUserLogin,
			UserDisplayName: event.FromBroadcasterUserName,
			Viewers:         int64(event.Viewers),
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
