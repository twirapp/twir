package twitch

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/libs/twitch"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	rewardsCacheKey      = "cache:twir:twitch:rewards:"
	rewardsCacheDuration = 6 * time.Hour
)

func BuildRewardsCacheKeyForId(twitchPlatformID string) string {
	return rewardsCacheKey + twitchPlatformID
}

func (c *CachedTwitchClient) GetChannelRewards(
	ctx context.Context,
	twitchUserID string,
	twitchPlatformID string,
) (
	[]helix.ChannelCustomReward,
	error,
) {
	if twitchUserID == "" || twitchPlatformID == "" {
		return nil, nil
	}

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("twitchUserID", twitchUserID),
		attribute.String("twitchPlatformID", twitchPlatformID),
	)

	if bytes, _ := c.redis.Get(ctx, BuildRewardsCacheKeyForId(twitchPlatformID)).Bytes(); len(bytes) > 0 {
		var rewards []helix.ChannelCustomReward
		if err := json.Unmarshal(bytes, &rewards); err != nil {
			return nil, err
		}

		return rewards, nil
	}

	twitchClient, err := twitch.NewUserClientWithContext(ctx, twitchUserID, c.config, c.twirBus)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create twitch client for broadcaster #%s: %w", twitchPlatformID, err,
		)
	}

	rewards, err := twitchClient.GetCustomRewards(
		&helix.GetCustomRewardsParams{
			BroadcasterID: twitchPlatformID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get rewards for broadcaster #%s: %w", twitchPlatformID, err,
		)
	}
	if rewards.ErrorMessage != "" {
		if rewards.StatusCode == http.StatusForbidden {
			return []helix.ChannelCustomReward{}, nil
		}

		return nil, fmt.Errorf(
			"failed to get rewards for broadcaster #%s: %s", twitchPlatformID, rewards.ErrorMessage,
		)
	}

	list := rewards.Data.ChannelCustomRewards

	for i, reward := range list {
		if reward.Image.Url1x == "" {
			list[i].Image.Url1x = reward.DefaultImage.Url1x
			list[i].Image.Url2x = reward.DefaultImage.Url2x
			list[i].Image.Url4x = reward.DefaultImage.Url4x
		}
	}

	rewardsBytes, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}

	if err := c.redis.Set(
		ctx,
		BuildRewardsCacheKeyForId(twitchPlatformID),
		rewardsBytes,
		rewardsCacheDuration,
	).Err(); err != nil {
		return nil, err
	}

	return rewards.Data.ChannelCustomRewards, nil
}
