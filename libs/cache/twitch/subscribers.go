package twitch

import (
	"context"
	"fmt"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/libs/twitch"
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
	twitchUserID string,
	twitchPlatformID string,
) (
	int,
	error,
) {
	if twitchUserID == "" || twitchPlatformID == "" {
		return 0, nil
	}

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("twitchUserID", twitchUserID),
		attribute.String("twitchPlatformID", twitchPlatformID),
	)

	if subscribers, err := c.redis.Get(
		ctx,
		buildChannelSubscribersCountCacheKeyForId(twitchPlatformID),
	).Int(); err == nil {
		return subscribers, nil
	}

	twitchClient, err := twitch.NewUserClient(twitchUserID, c.config, c.twirBus)
	if err != nil {
		return 0, fmt.Errorf("failed to create twitch client: %w", err)
	}

	subscribersReq, err := twitchClient.GetSubscriptions(
		&helix.SubscriptionsParams{
			BroadcasterID: twitchPlatformID,
		},
	)
	if err != nil {
		return 0, err
	}
	if subscribersReq.ErrorMessage != "" {
		return 0, fmt.Errorf("cannot get channels subscribers: %s", subscribersReq.ErrorMessage)
	}

	if err := c.redis.Set(
		ctx,
		buildChannelSubscribersCountCacheKeyForId(twitchPlatformID),
		subscribersReq.Data.Total,
		channelSubscribersCountCacheDuration,
	).Err(); err != nil {
		return 0, err
	}

	return subscribersReq.Data.Total, nil
}
