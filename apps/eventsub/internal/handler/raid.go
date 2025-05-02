package handler

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/events"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelRaid(
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

	if err := c.gorm.Create(
		&model.ChannelsEventsListItem{
			ID:        uuid.New().String(),
			ChannelID: event.ToBroadcasterUserID,
			UserID:    event.FromBroadcasterUserID,
			Type:      model.ChannelEventListItemTypeRaided,
			CreatedAt: time.Now().UTC(),
			Data: &model.ChannelsEventsListItemData{
				RaidedViewersCount:    strconv.Itoa(event.Viewers),
				RaidedFromUserName:    event.FromBroadcasterUserLogin,
				RaidedFromDisplayName: event.FromBroadcasterUserName,
			},
		},
	).Error; err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	if err := c.twirBus.Events.Raided.Publish(
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
