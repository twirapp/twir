package handler

import (
	"context"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/twirapp/twir/libs/cache/twitch"
)

func (c *Handler) flushChannelPointsRewardsCache(ctx context.Context, channelID string) error {
	return c.redisClient.Del(ctx, twitch.BuildRewardsCacheKeyForId(channelID)).Err()
}

func (c *Handler) HandleChannelPointsRewardAdd(
	ctx context.Context,
	event eventsub.ChannelPointsCustomRewardAddEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"Channel points reward added",
		slog.String("reward", event.Title),
		slog.String("channel_id", event.BroadcasterUserId),
	)

	if err := c.flushChannelPointsRewardsCache(ctx, event.BroadcasterUserId); err != nil {
		c.logger.Error("failed to flush channel points rewards cache", slog.Any("err", err))
	}
}

func (c *Handler) HandleChannelPointsRewardUpdate(
	ctx context.Context,
	event eventsub.ChannelPointsCustomRewardUpdateEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"Channel points reward updated",
		slog.String("reward", event.Title),
		slog.String("channel_id", event.BroadcasterUserId),
	)

	if err := c.flushChannelPointsRewardsCache(ctx, event.BroadcasterUserId); err != nil {
		c.logger.Error("failed to flush channel points rewards cache", slog.Any("err", err))
	}
}

func (c *Handler) HandleChannelPointsRewardRemove(
	ctx context.Context,
	event eventsub.ChannelPointsCustomRewardRemoveEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"Channel points reward removed",
		slog.String("reward", event.Title),
		slog.String("channel_id", event.BroadcasterUserId),
	)

	if err := c.flushChannelPointsRewardsCache(ctx, event.BroadcasterUserId); err != nil {
		c.logger.Error("failed to flush channel points rewards cache", slog.Any("err", err))
	}
}
