package handler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"

	"github.com/twirapp/twir/libs/grpc/events"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelFollow(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelFollow,
) {
	redisKey := fmt.Sprintf("follows-cache:%s:%s", event.BroadcasterUserID, event.UserID)
	key, _ := c.redisClient.Get(context.Background(), redisKey).Result()

	if key != "" {
		return
	}

	c.logger.Info(
		"channel follow",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("userId", event.UserID),
		slog.String("userName", event.UserLogin),
	)

	c.gorm.Create(
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

	c.eventsGrpc.Follow(
		ctx,
		&events.FollowMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			UserName:        event.UserLogin,
			UserDisplayName: event.UserName,
			UserId:          event.UserID,
		},
	)

	c.redisClient.Set(ctx, redisKey, redisKey, 24*7*time.Hour)
}
