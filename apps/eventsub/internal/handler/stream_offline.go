package handler

import (
	"context"
	"encoding/json"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/pubsub"
	"github.com/twirapp/twir/libs/grpc/events"
	"go.uber.org/zap"
)

func (c *Handler) handleStreamOffline(
	h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventStreamOffline,
) {
	defer zap.S().Infow(
		"stream offline", ""+
			"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
	)

	c.services.Grpc.Events.StreamOffline(
		context.Background(), &events.StreamOfflineMessage{
			BaseInfo: &events.BaseInfo{ChannelId: event.BroadcasterUserID},
		},
	)

	bytes, err := json.Marshal(
		&pubsub.StreamOfflineMessage{
			ChannelID: event.BroadcasterUserID,
		},
	)
	if err != nil {
		zap.S().Error(err)
		return
	}

	err = c.services.Gorm.Where(
		`"userId" = ?`,
		event.BroadcasterUserID,
	).Delete(&model.ChannelsStreams{}).Error
	if err != nil {
		zap.S().Error(err)
	}

	c.services.PubSub.Publish("stream.offline", bytes)
}
