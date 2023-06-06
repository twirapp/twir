package handler

import (
	"context"
	"fmt"
	"github.com/dnsge/twitch-eventsub-bindings"
	"github.com/satont/tsuwari/libs/grpc/generated/events"
	"go.uber.org/zap"
	"time"
)

func (c *Handler) handleChannelFollow(h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventChannelFollow) {
	redisKey := fmt.Sprintf("follows-cache:%s:%s", event.BroadcasterUserID, event.UserID)
	key, _ := c.services.Redis.Get(context.Background(), redisKey).Result()

	if key != "" {
		return
	}

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

	c.services.Redis.Set(context.Background(), redisKey, redisKey, 24*7*time.Hour)
}
