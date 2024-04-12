package twitch

import (
	"context"
	"fmt"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/libs/twitch"
)

const channelFollowersCountCacheKey = "cache:twir:twitch:followersCount:"
const channelFollowersCountCacheDuration = 10 * time.Minute

func buildChannelFollowersCountCacheKeyForId(userId string) string {
	return channelFollowersCountCacheKey + userId
}

func (c *CachedTwitchClient) GetChannelFollowersCountByChannelId(
	ctx context.Context,
	channelId string,
) (
	int,
	error,
) {
	if channelId == "" {
		return 0, nil
	}

	if followers, err := c.redis.Get(
		ctx,
		buildChannelFollowersCountCacheKeyForId(channelId),
	).Int(); err == nil {
		return followers, nil
	}

	twitchClient, err := twitch.NewUserClient(channelId, c.config, c.tokensClient)
	if err != nil {
		return 0, fmt.Errorf("failed to create twitch client: %w", err)
	}

	followsReq, err := twitchClient.GetChannelFollows(
		&helix.GetChannelFollowsParams{
			BroadcasterID: channelId,
		},
	)
	if err != nil {
		return 0, err
	}
	if followsReq.ErrorMessage != "" {
		return 0, fmt.Errorf("cannot get channels followers: %s", followsReq.ErrorMessage)
	}

	if err := c.redis.Set(
		ctx,
		buildChannelFollowersCountCacheKeyForId(channelId),
		followsReq.Data.Total,
		channelFollowersCountCacheDuration,
	).Err(); err != nil {
		return 0, err
	}

	return followsReq.Data.Total, nil
}
