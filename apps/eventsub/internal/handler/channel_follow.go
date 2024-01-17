package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	"github.com/twirapp/twir/libs/grpc/events"
	"go.uber.org/zap"
)

func (c *Handler) handleChannelFollow(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelFollow,
) {
	redisKey := fmt.Sprintf("follows-cache:%s:%s", event.BroadcasterUserID, event.UserID)
	key, _ := c.services.Redis.Get(context.Background(), redisKey).Result()

	if key != "" {
		return
	}

	zap.S().Infow(
		"channel follow",
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
		"userId", event.UserID,
		"userName", event.UserLogin,
	)

	c.services.Gorm.Create(
		&model.ChannelsEventsListItem{
			ID:        uuid.New().String(),
			ChannelID: event.BroadcasterUserID,
			UserID:    event.UserID,
			Type:      model.ChannelEventListItemTypeFollow,
			CreatedAt: time.Now().UTC(),
			Data: &model.ChannelsEventsListItemData{
				FollowUserName:        event.UserLogin,
				FollowUserDisplayName: event.UserName,
			},
		},
	)

	ctx := context.Background()

	c.services.Grpc.Events.Follow(
		ctx,
		&events.FollowMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			UserName:        event.UserLogin,
			UserDisplayName: event.UserName,
			UserId:          event.UserID,
		},
	)

	c.services.Redis.Set(ctx, redisKey, redisKey, 24*7*time.Hour)
}
