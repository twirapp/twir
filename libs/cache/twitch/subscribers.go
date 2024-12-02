package twitch

import (
	"context"
	"fmt"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/libs/twitch"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const channelSubscribersCountCacheKey = "cache:twir:twitch:subscribersCount:"
const channelSubscribersCountCacheDuration = 10 * time.Minute

func buildChannelSubscribersCountCacheKeyForId(userId string) string {
	return channelSubscribersCountCacheKey + userId
}

func (c *CachedTwitchClient) GetChannelSubscribersCountByChannelId(
	ctx context.Context,
	channelId string,
) (
	int,
	error,
) {
	if channelId == "" {
		return 0, nil
	}

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("channelId", channelId),
	)

	if subscribers, err := c.redis.Get(
		ctx,
		buildChannelSubscribersCountCacheKeyForId(channelId),
	).Int(); err == nil {
		return subscribers, nil
	}

	twitchClient, err := twitch.NewUserClient(channelId, c.config, c.tokensClient)
	if err != nil {
		return 0, fmt.Errorf("failed to create twitch client: %w", err)
	}

	subscribersReq, err := twitchClient.GetSubscriptions(
		&helix.SubscriptionsParams{
			BroadcasterID: channelId,
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
		buildChannelSubscribersCountCacheKeyForId(channelId),
		subscribersReq.Data.Total,
		channelSubscribersCountCacheDuration,
	).Err(); err != nil {
		return 0, err
	}

	return subscribersReq.Data.Total, nil
}
