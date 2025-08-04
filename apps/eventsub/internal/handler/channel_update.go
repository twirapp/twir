package handler

import (
	"context"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/bus-core/events"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/redis_keys"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
)

func (c *Handler) HandleChannelUpdate(
	ctx context.Context,
	event eventsub.ChannelUpdateEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	if err := c.redisClient.Del(
		ctx,
		redis_keys.StreamByChannelID(event.BroadcasterUserId),
	).Err(); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
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
			BaseInfo:    events.BaseInfo{ChannelID: event.BroadcasterUserId},
			NewTitle:    event.Title,
			NewCategory: event.CategoryName,
		},
	)

	if err := c.channelsInfoHistoryRepo.Create(
		ctx,
		channelsinfohistory.CreateInput{
			ChannelID: event.BroadcasterUserId,
			Title:     event.Title,
			Category:  event.CategoryName,
		},
	); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	err := c.gorm.
		WithContext(ctx).
		Model(&model.ChannelsStreams{}).
		Where(`"userId" = ?`, event.BroadcasterUserId).
		Updates(
			map[string]any{
				"title":      event.Title,
				`"gameName"`: event.CategoryName,
				`"gameId"`:   event.CategoryId,
			},
		).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
