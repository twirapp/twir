package handler

import (
	"context"
	"github.com/dnsge/twitch-eventsub-bindings"
	"github.com/satont/tsuwari/libs/grpc/generated/events"
	"go.uber.org/zap"
)

func (c *Handler) handleChannelFollow(h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelFollow) {
	zap.S().Infow("channel follow",
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
		"userId", event.UserID,
		"userName", event.UserLogin,
	)

	c.services.Grpc.Events.Follow(context.Background(), &events.FollowMessage{
		BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
		UserName:        event.UserLogin,
		UserDisplayName: event.UserName,
		UserId:          event.UserID,
	})
}
