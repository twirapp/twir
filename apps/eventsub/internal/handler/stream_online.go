package handler

import (
	"context"
	"encoding/json"
	"github.com/dnsge/twitch-eventsub-bindings"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/libs/grpc/generated/events"
	"github.com/satont/tsuwari/libs/pubsub"
	"github.com/satont/tsuwari/libs/twitch"
	"go.uber.org/zap"
)

func (c *handler) handleStreamOnline(h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventStreamOnline) {
	twitchClient, err := twitch.NewAppClient(*c.services.Config, c.services.Grpc.Tokens)
	if err != nil {
		zap.S().Error(err)
		return
	}

	streamsReq, err := twitchClient.GetStreams(&helix.StreamsParams{
		UserIDs: []string{event.BroadcasterUserID},
	})

	if err != nil || streamsReq.ErrorMessage != "" {
		zap.S().Error(err, streamsReq.ErrorMessage)
		return
	}

	if len(streamsReq.Data.Streams) == 0 {
		return
	}

	c.services.Grpc.Events.StreamOnline(context.Background(), &events.StreamOnlineMessage{
		BaseInfo: &events.BaseInfo{ChannelId: event.BroadcasterUserID},
		Title:    streamsReq.Data.Streams[0].Title,
		Category: streamsReq.Data.Streams[0].GameName,
	})

	bytes, err := json.Marshal(&pubsub.StreamOnlineMessage{
		ChannelID: event.BroadcasterUserID,
		StreamID:  event.ID,
	})
	if err != nil {
		zap.S().Error(err)
		return
	}

	c.services.PubSub.Publish("stream.online", bytes)
}
