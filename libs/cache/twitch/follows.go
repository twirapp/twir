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

const channelFollowersCountCacheKey = "cache:twir:twitch:followersCount:"
const channelFollowersCountCacheDuration = 10 * time.Minute

func buildChannelFollowersCountCacheKeyForId(twitchPlatformID string) string {
	return channelFollowersCountCacheKey + twitchPlatformID
}

func (c *CachedTwitchClient) GetChannelFollowersCountByChannelId(
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
		attribute.String("twitch.twitchUserID", twitchUserID),
		attribute.String("twitch.twitchPlatformID", twitchPlatformID),
	)

	if followers, err := c.redis.Get(
		ctx,
		buildChannelFollowersCountCacheKeyForId(twitchPlatformID),
	).Int(); err == nil {
		return followers, nil
	}

	twitchClient, err := twitch.NewUserClient(twitchUserID, c.config, c.twirBus)
	if err != nil {
		return 0, fmt.Errorf("failed to create twitch client: %w", err)
	}

	followsReq, err := twitchClient.GetChannelFollows(
		&helix.GetChannelFollowsParams{
			BroadcasterID: twitchPlatformID,
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
		buildChannelFollowersCountCacheKeyForId(twitchPlatformID),
		followsReq.Data.Total,
		channelFollowersCountCacheDuration,
	).Err(); err != nil {
		return 0, err
	}

	return followsReq.Data.Total, nil
}

// GetUserFollowDuration returns the duration a user has been following a channel
// Returns nil if the user is not following the channel
func (c *CachedTwitchClient) GetUserFollowDuration(
	ctx context.Context,
	twitchUserID string,
	followerPlatformID string,
	channelPlatformID string,
) (*time.Duration, error) {
	if twitchUserID == "" || followerPlatformID == "" || channelPlatformID == "" {
		return nil, nil
	}

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("twitch.twitchUserID", twitchUserID),
		attribute.String("twitch.followerPlatformID", followerPlatformID),
		attribute.String("twitch.channelPlatformID", channelPlatformID),
	)

	twitchClient, err := twitch.NewUserClient(twitchUserID, c.config, c.twirBus)
	if err != nil {
		return nil, fmt.Errorf("failed to create twitch client: %w", err)
	}

	// Use GetUsersFollows to check if a specific user is following
	followsReq, err := twitchClient.GetUsersFollows(
		&helix.UsersFollowsParams{
			FromID: followerPlatformID,
			ToID:   channelPlatformID,
		},
	)
	if err != nil {
		return nil, err
	}
	if followsReq.ErrorMessage != "" {
		return nil, fmt.Errorf("cannot get user follow: %s", followsReq.ErrorMessage)
	}

	// User is not following
	if len(followsReq.Data.Follows) == 0 {
		return nil, nil
	}

	followedAt := followsReq.Data.Follows[0].FollowedAt
	duration := time.Since(followedAt)

	return &duration, nil
}
