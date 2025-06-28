package handler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/twirapp/twir/libs/bus-core/events"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"

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

	if err := c.eventsListRepository.Create(
		context.TODO(),
		channelseventslist.CreateInput{
			ChannelID: event.BroadcasterUserID,
			UserID:    &event.UserID,
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
