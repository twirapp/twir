package handler

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) flushChannelPointsRewardsCache(ctx context.Context, channelID string) error {
	return c.redisClient.Del(ctx, twitch.BuildRewardsCacheKeyForId(channelID)).Err()
}

func (c *Handler) handleChannelPointsRewardAdd(
	ctx context.Context,
	_ *esb.ResponseHeaders,
	event *esb.EventChannelPointsRewardAdd,
) {
	c.logger.Info(
		"Channel points reward added",
		slog.String("reward", event.Title),
		slog.String("channel_id", event.BroadcasterUserID),
	)

	if err := c.flushChannelPointsRewardsCache(ctx, event.BroadcasterUserID); err != nil {
		c.logger.Error("failed to flush channel points rewards cache", slog.Any("err", err))
	}
}

func (c *Handler) handleChannelPointsRewardUpdate(
	ctx context.Context,
	_ *esb.ResponseHeaders,
	event *esb.EventChannelPointsRewardUpdate,
) {
	c.logger.Info(
		"Channel points reward updated",
		slog.String("reward", event.Title),
		slog.String("channel_id", event.BroadcasterUserID),
	)

	if err := c.flushChannelPointsRewardsCache(ctx, event.BroadcasterUserID); err != nil {
		c.logger.Error("failed to flush channel points rewards cache", slog.Any("err", err))
	}
}

func (c *Handler) handleChannelPointsRewardRemove(
	ctx context.Context,
	_ *esb.ResponseHeaders,
	event *esb.EventChannelPointsRewardRemove,
) {
	c.logger.Info(
		"Channel points reward removed",
		slog.String("reward", event.Title),
		slog.String("channel_id", event.BroadcasterUserID),
	)

	if err := c.flushChannelPointsRewardsCache(ctx, event.BroadcasterUserID); err != nil {
		c.logger.Error("failed to flush channel points rewards cache", slog.Any("err", err))
	}
}
