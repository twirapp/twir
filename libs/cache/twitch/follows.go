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

const channelFollowersCountCacheKey = "cache:twir:twitch:followersCount:"
const channelFollowersCountCacheDuration = 10 * time.Minute

func buildChannelFollowersCountCacheKeyForId(twitchPlatformID string) string {
	return channelFollowersCountCacheKey + twitchPlatformID
}

func (c *CachedTwitchClient) GetChannelFollowersCountByChannelId(
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
		attribute.String("twitch.twitchUserID", twitchUserID.String()),
		attribute.String("twitch.twitchPlatformID", twitchPlatformID),
	)

	if followers, err := c.redis.Get(
		ctx,
		buildChannelFollowersCountCacheKeyForId(twitchPlatformID),
	).Int(); err == nil {
		return followers, nil
	}

	twitchClient, err := c.createUserClient(ctx, twitchUserID)
	if err != nil {
		return 0, fmt.Errorf("failed to create twitch client: %w", err)
	}

	first := 100
	pager, err := twitchClient.Channels.GetChannelFollowersPager(
		helix.GetChannelFollowersRequest{
			BroadcasterID: twitchPlatformID,
			First:         &first,
		},
		helix.WithPageLimit(10000),
	)
	if err != nil {
		return 0, err
	}

	followers := 0
	for pager.Next(ctx) {
		followers += len(pager.Page().Data)
	}
	if err := pager.Err(); err != nil {
		return 0, err
	}

	if err := c.redis.Set(
		ctx,
		buildChannelFollowersCountCacheKeyForId(twitchPlatformID),
		followers,
		channelFollowersCountCacheDuration,
	).Err(); err != nil {
		return 0, err
	}

	return followers, nil
}

// GetUserFollowDuration returns the duration a user has been following a channel
// Returns nil if the user is not following the channel
func (c *CachedTwitchClient) GetUserFollowDuration(
	ctx context.Context,
	twitchUserID uuid.UUID,
	followerPlatformID string,
	channelPlatformID string,
) (*time.Duration, error) {
	if twitchUserID == uuid.Nil || followerPlatformID == "" || channelPlatformID == "" {
		return nil, nil
	}

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("twitch.twitchUserID", twitchUserID.String()),
		attribute.String("twitch.followerPlatformID", followerPlatformID),
		attribute.String("twitch.channelPlatformID", channelPlatformID),
	)

	twitchClient, err := c.createUserClient(ctx, twitchUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create twitch client: %w", err)
	}

	followsReq, err := twitchClient.Channels.GetFollowedChannels(
		ctx, helix.GetFollowedChannelsRequest{
			UserID:        followerPlatformID,
			BroadcasterID: &channelPlatformID,
		},
	)
	if err != nil {
		return nil, err
	}

	// User is not following
	if len(followsReq.Data) == 0 {
		return nil, nil
	}

	followedAt := followsReq.Data[0].FollowedAt.Time
	duration := time.Since(followedAt)

	return &duration, nil
}
