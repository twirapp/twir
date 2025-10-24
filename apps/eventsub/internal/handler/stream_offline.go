package handler

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/bus-core/twitch"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/redis_keys"
	"github.com/twirapp/twitchy/eventsub"
)

func (c *Handler) HandleStreamOffline(
	ctx context.Context,
	event eventsub.StreamOfflineEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"stream offline",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("channelName", event.BroadcasterUserLogin),
	)

	if err := c.redisClient.Del(
		ctx,
		redis_keys.StreamByChannelID(event.BroadcasterUserId),
	).Err(); err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}

	dbStream := model.ChannelsStreams{}
	if err := c.gorm.WithContext(ctx).Where(
		`"userId" = ?`,
		event.BroadcasterUserId,
	).First(&dbStream).Error; err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
		return
	}

	c.twirBus.Channel.StreamOffline.Publish(
		ctx,
		twitch.StreamOfflineMessage{
			ChannelID: event.BroadcasterUserId,
			StartedAt: dbStream.StartedAt,
		},
	)

	err := c.gorm.WithContext(ctx).Where(
		`"userId" = ?`,
		event.BroadcasterUserId,
	).Delete(&model.ChannelsStreams{}).Error
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
