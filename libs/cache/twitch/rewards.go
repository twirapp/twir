package twitch

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/helix"
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
	twitchUserID uuid.UUID,
	twitchPlatformID string,
) (
	[]helix.CustomReward,
	error,
) {
	if twitchUserID == uuid.Nil || twitchPlatformID == "" {
		return nil, nil
	}

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("twitchUserID", twitchUserID.String()),
		attribute.String("twitchPlatformID", twitchPlatformID),
	)

	if bytes, _ := c.redis.Get(ctx, BuildRewardsCacheKeyForId(twitchPlatformID)).Bytes(); len(bytes) > 0 {
		var rewards []helix.CustomReward
		if err := json.Unmarshal(bytes, &rewards); err != nil {
			return nil, err
		}

		return rewards, nil
	}

	twitchClient, err := c.createUserClient(ctx, twitchUserID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create twitch client for broadcaster #%s: %w", twitchPlatformID, err,
		)
	}

	rewards, err := twitchClient.ChannelPoints.GetCustomReward(
		ctx, helix.GetCustomRewardRequest{
			BroadcasterID: twitchPlatformID,
		},
	)
	if err != nil {
		var authErr *helix.AuthError
		if errors.As(err, &authErr) && authErr.StatusCode() == http.StatusForbidden && strings.Contains(err.Error(), "The broadcaster must have partner or affiliate status.") {
			return []helix.CustomReward{}, nil
		}

		return nil, fmt.Errorf(
			"failed to get rewards for broadcaster #%s: %w", twitchPlatformID, err,
		)
	}

	list := rewards.Data

	for i, reward := range list {
		if reward.Image == nil || reward.Image.URL1x == "" {
			list[i].Image = &helix.CustomRewardImage{
				URL1x: reward.DefaultImage.URL1x,
				URL2x: reward.DefaultImage.URL2x,
				URL4x: reward.DefaultImage.URL4x,
			}
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

	return list, nil
}
