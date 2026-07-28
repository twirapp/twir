package handler

import (
	"context"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/events"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/redis_keys"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
)

func (c *Handler) HandleChannelUpdate(
	ctx context.Context,
	event eventsub.ChannelUpdateEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	channel, err := c.channelService.GetChannelByPlatformChannelID(
		ctx,
		platformentity.PlatformTwitch,
		event.BroadcasterUserId,
	)
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
		return
	}
	channelID := channel.ID.String()

	if err := c.redisClient.Del(
		ctx,
		redis_keys.StreamByChannelID(channelID),
	).Err(); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}

	c.logger.Info(
		"channel update",
		slog.String("title", event.Title),
		slog.String("category", event.CategoryName),
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("channelName", event.BroadcasterUserLogin),
	)

	c.twirBus.Events.TitleOrCategoryChanged.Publish(
		ctx,
		events.TitleOrCategoryChangedMessage{
			BaseInfo:    events.BaseInfo{ChannelPlatformID: event.BroadcasterUserId, Platform: platformentity.PlatformTwitch},
			NewTitle:    event.Title,
			NewCategory: event.CategoryName,
		},
	)

	if err := c.channelsInfoHistoryRepo.Create(
		ctx,
		channelsinfohistory.CreateInput{
			ChannelID: channelID,
			Platform:  platformentity.PlatformTwitch,
			Title:     event.Title,
			Category:  event.CategoryName,
		},
	); err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}

	err = c.streamsrepository.Update(
		ctx,
		channel.ID,
		platformentity.PlatformTwitch,
		streamsrepository.UpdateInput{
			Title:    &event.Title,
			GameName: &event.CategoryName,
			GameId:   &event.CategoryId,
		},
	)
	if err != nil {
		c.logger.Error(err.Error(), logger.Error(err))
	}
}
