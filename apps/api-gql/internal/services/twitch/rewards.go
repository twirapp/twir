package twitch

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/helix"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type CustomRewardsResult struct {
	Rewards              []helix.CustomReward
	IsPartnerOrAffiliate bool
}

func (c *Service) GetRewardsByChannelID(
	ctx context.Context,
	channelID string,
) (CustomRewardsResult, error) {
	parsedID, err := uuid.Parse(channelID)
	if err != nil {
		return CustomRewardsResult{}, fmt.Errorf("invalid channel id: %w", err)
	}

	channel, err := c.channelService.GetChannelByID(ctx, parsedID)
	if err != nil {
		return CustomRewardsResult{}, fmt.Errorf("get channel: %w", err)
	}
	if channel.IsNil() {
		return CustomRewardsResult{}, nil
	}

	twitchBinding, found := channel.Binding(platformentity.PlatformTwitch)
	if !found || twitchBinding.UserID == uuid.Nil {
		return CustomRewardsResult{}, nil
	}

	twitchClient, err := c.createUserClient(ctx, twitchBinding.UserID)
	if err != nil {
		return CustomRewardsResult{}, fmt.Errorf("failed to create twitch client: %w", err)
	}

	rewards, err := twitchClient.ChannelPoints.GetCustomReward(
		ctx, helix.GetCustomRewardRequest{
			BroadcasterID: twitchBinding.PlatformChannelID,
		},
	)
	if err != nil {
		var authErr *helix.AuthError
		if errors.As(err, &authErr) && authErr.StatusCode() == 403 && strings.Contains(err.Error(), "The broadcaster must have partner or affiliate status.") {
			return CustomRewardsResult{IsPartnerOrAffiliate: false}, nil
		}

		return CustomRewardsResult{}, fmt.Errorf("cannot get custom rewards: %w", err)
	}

	return CustomRewardsResult{
		Rewards:              rewards.Data,
		IsPartnerOrAffiliate: true,
	}, nil
}
