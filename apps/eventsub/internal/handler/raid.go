package handler

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/twirapp/twir/libs/bus-core/events"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
	"github.com/twirapp/twitchy/eventsub"
)

func (c *Handler) HandleChannelRaid(
	ctx context.Context,
	event eventsub.ChannelRaidEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"channel raid",
		slog.String("channelId", event.ToBroadcasterUserId),
		slog.String("channelName", event.ToBroadcasterUserName),
		slog.String("userId", event.FromBroadcasterUserId),
		slog.String("userName", event.FromBroadcasterUserLogin),
		slog.Int("viewers", event.Viewers),
	)

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.ToBroadcasterUserId,
			UserID:    &event.FromBroadcasterUserId,
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
				ChannelName: event.ToBroadcasterUserId,
				ChannelID:   event.ToBroadcasterUserId,
			},
			UserID:          event.FromBroadcasterUserId,
			UserName:        event.FromBroadcasterUserLogin,
			UserDisplayName: event.FromBroadcasterUserName,
			Viewers:         int64(event.Viewers),
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
