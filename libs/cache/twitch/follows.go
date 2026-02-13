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

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("twitch.channelId", channelId),
	)

	if followers, err := c.redis.Get(
		ctx,
		buildChannelFollowersCountCacheKeyForId(channelId),
	).Int(); err == nil {
		return followers, nil
	}

	twitchClient, err := twitch.NewUserClient(channelId, c.config, c.twirBus)
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

// GetUserFollowDuration returns the duration a user has been following a channel
// Returns nil if the user is not following the channel
func (c *CachedTwitchClient) GetUserFollowDuration(
	ctx context.Context,
	userID string,
	channelID string,
) (*time.Duration, error) {
	if userID == "" || channelID == "" {
		return nil, nil
	}

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("twitch.userId", userID),
		attribute.String("twitch.channelId", channelID),
	)

	twitchClient, err := twitch.NewUserClient(channelID, c.config, c.twirBus)
	if err != nil {
		return nil, fmt.Errorf("failed to create twitch client: %w", err)
	}

	// Use GetUsersFollows to check if a specific user is following
	followsReq, err := twitchClient.GetUsersFollows(
		&helix.UsersFollowsParams{
			FromID: userID,
			ToID:   channelID,
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
