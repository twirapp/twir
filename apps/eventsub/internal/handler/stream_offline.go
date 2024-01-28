package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/pubsub"
	"github.com/twirapp/twir/libs/grpc/events"
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

	bytes, err := json.Marshal(
		&pubsub.StreamOfflineMessage{
			ChannelID: event.BroadcasterUserID,
		},
	)
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}

	err = c.gorm.Where(
		`"userId" = ?`,
		event.BroadcasterUserID,
	).Delete(&model.ChannelsStreams{}).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	c.pubSub.Client.Publish("stream.offline", bytes)
}
