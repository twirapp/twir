package twitch

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/helix"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const channelSubscribersCountCacheKey = "cache:twir:twitch:subscribersCount:"
const channelSubscribersCountCacheDuration = 10 * time.Minute

func buildChannelSubscribersCountCacheKeyForId(twitchPlatformID string) string {
	return channelSubscribersCountCacheKey + twitchPlatformID
}

func (c *CachedTwitchClient) GetChannelSubscribersCountByChannelId(
	ctx context.Context,
	twitchUserID uuid.UUID,
	twitchPlatformID string,
) (
	int,
	error,
) {
	if twitchUserID == uuid.Nil || twitchPlatformID == "" {
		return 0, nil
	}

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("twitchUserID", twitchUserID.String()),
		attribute.String("twitchPlatformID", twitchPlatformID),
	)

	if subscribers, err := c.redis.Get(
		ctx,
		buildChannelSubscribersCountCacheKeyForId(twitchPlatformID),
	).Int(); err == nil {
		return subscribers, nil
	}

	twitchClient, err := c.createUserClient(ctx, twitchUserID)
	if err != nil {
		return 0, fmt.Errorf("failed to create twitch client: %w", err)
	}

	subscribersReq, err := twitchClient.Subscriptions.GetBroadcasterSubscriptions(
		ctx, helix.GetBroadcasterSubscriptionsRequest{
			BroadcasterID: twitchPlatformID,
		},
	)
	if err != nil {
		return 0, err
	}

	subscribers := 0
	if subscribersReq.Data.Total != nil {
		subscribers = *subscribersReq.Data.Total
	}

	if err := c.redis.Set(
		ctx,
		buildChannelSubscribersCountCacheKeyForId(twitchPlatformID),
		subscribers,
		channelSubscribersCountCacheDuration,
	).Err(); err != nil {
		return 0, err
	}

	return subscribers, nil
}
