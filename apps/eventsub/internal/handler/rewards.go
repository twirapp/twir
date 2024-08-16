package handler

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twitch-eventsub-framework/esb"
)

func (c *Handler) flushChannelPointsRewardsCache(channelID string) error {
	return c.redisClient.Del(context.TODO(), twitch.BuildRewardsCacheKeyForId(channelID)).Err()
}

func (c *Handler) handleChannelPointsRewardAdd(
	_ *esb.ResponseHeaders,
	event *esb.EventChannelPointsRewardAdd,
) {
	c.logger.Info(
		"Channel points reward added",
		slog.String("reward", event.Title),
		slog.String("channel_id", event.BroadcasterUserID),
	)

	if err := c.flushChannelPointsRewardsCache(event.BroadcasterUserID); err != nil {
		c.logger.Error("failed to flush channel points rewards cache", slog.Any("err", err))
	}
}

func (c *Handler) handleChannelPointsRewardUpdate(
	_ *esb.ResponseHeaders,
	event *esb.EventChannelPointsRewardUpdate,
) {
	c.logger.Info(
		"Channel points reward updated",
		slog.String("reward", event.Title),
		slog.String("channel_id", event.BroadcasterUserID),
	)

	if err := c.flushChannelPointsRewardsCache(event.BroadcasterUserID); err != nil {
		c.logger.Error("failed to flush channel points rewards cache", slog.Any("err", err))
	}
}

func (c *Handler) handleChannelPointsRewardRemove(
	_ *esb.ResponseHeaders,
	event *esb.EventChannelPointsRewardRemove,
) {
	c.logger.Info(
		"Channel points reward removed",
		slog.String("reward", event.Title),
		slog.String("channel_id", event.BroadcasterUserID),
	)

	if err := c.flushChannelPointsRewardsCache(event.BroadcasterUserID); err != nil {
		c.logger.Error("failed to flush channel points rewards cache", slog.Any("err", err))
	}
}
