package twitch

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/libs/twitch"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const rewardsCacheKey = "cache:twir:twitch:rewards:"
const rewardsCacheDuration = 6 * time.Hour

func BuildRewardsCacheKeyForId(channelId string) string {
	return rewardsCacheKey + channelId
}

func (c *CachedTwitchClient) GetChannelRewards(
	ctx context.Context,
	channelID string,
) (
	[]helix.ChannelCustomReward,
	error,
) {
	if channelID == "" {
		return nil, nil
	}

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("channelID", channelID),
	)

	if bytes, _ := c.redis.Get(ctx, BuildRewardsCacheKeyForId(channelID)).Bytes(); len(bytes) > 0 {
		var rewards []helix.ChannelCustomReward
		if err := json.Unmarshal(bytes, &rewards); err != nil {
			return nil, err
		}

		return rewards, nil
	}

	twitchClient, err := twitch.NewUserClientWithContext(ctx, channelID, c.config, c.twirBus)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create twitch client for broadcaster #%s: %w", channelID, err,
		)
	}

	rewards, err := twitchClient.GetCustomRewards(
		&helix.GetCustomRewardsParams{
			BroadcasterID: channelID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get rewards for broadcaster #%s: %w", channelID, err,
		)
	}
	if rewards.ErrorMessage != "" {
		return nil, fmt.Errorf(
			"failed to get rewards for broadcaster #%s: %s", channelID, rewards.ErrorMessage,
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
		BuildRewardsCacheKeyForId(channelID),
		rewardsBytes,
		rewardsCacheDuration,
	).Err(); err != nil {
		return nil, err
	}

	return rewards.Data.ChannelCustomRewards, nil
}
