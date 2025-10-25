package handler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/events"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
)

func (c *Handler) HandleChannelFollow(
	ctx context.Context,
	event eventsub.ChannelFollowEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	redisKey := fmt.Sprintf("follows-cache:%s:%s", event.BroadcasterUserId, event.UserId)
	key, _ := c.redisClient.Get(ctx, redisKey).Result()

	if key != "" {
		return
	}

	c.redisClient.Set(ctx, redisKey, redisKey, 24*7*time.Hour)

	c.logger.Info(
		"channel follow",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("userId", event.UserId),
		slog.String("userName", event.UserLogin),
	)

	if err := c.eventsListRepository.Create(
		ctx,
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserId,
			UserID:    &event.UserId,
			Type:      model.ChannelEventListItemTypeFollow,
			Data: &model.ChannelsEventsListItemData{
				FollowUserName:        event.UserLogin,
				FollowUserDisplayName: event.UserName,
			},
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	c.twirBus.Events.Follow.Publish(
		ctx,
		events.FollowMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserName:        event.UserLogin,
			UserDisplayName: event.UserName,
			UserID:          event.UserId,
		},
	)
}
