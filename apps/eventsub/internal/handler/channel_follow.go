package handler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/events"

	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) handleChannelFollow(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelFollow,
) {
	ctx := context.Background()
	redisKey := fmt.Sprintf("follows-cache:%s:%s", event.BroadcasterUserID, event.UserID)
	key, _ := c.redisClient.Get(ctx, redisKey).Result()

	if key != "" {
		return
	}

	c.redisClient.Set(ctx, redisKey, redisKey, 24*7*time.Hour)

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

	c.twirBus.Events.Follow.Publish(
		events.FollowMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserName:        event.UserLogin,
			UserDisplayName: event.UserName,
			UserID:          event.UserID,
		},
	)
}
