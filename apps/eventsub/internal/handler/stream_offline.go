package handler

import (
	"context"
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/grpc/events"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleStreamOffline(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventStreamOffline,
) {
	c.logger.Info(
		"stream offline",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
	)

	_, err := c.eventsGrpc.StreamOffline(
		context.Background(),
		&events.StreamOfflineMessage{
			BaseInfo: &events.BaseInfo{ChannelId: event.BroadcasterUserID},
		},
	)
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	err = c.gorm.Where(
		`"userId" = ?`,
		event.BroadcasterUserID,
	).Delete(&model.ChannelsStreams{}).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	c.bus.Channel.StreamOffline.Publish(twitch.StreamOfflineMessage{ChannelID: event.BroadcasterUserID})
}
