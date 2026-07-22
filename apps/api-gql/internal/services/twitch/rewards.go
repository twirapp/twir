package twitch

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	apiChannelbinding "github.com/twirapp/twir/apps/api-gql/internal/channelbinding"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type CustomRewardsResult struct {
	Rewards              []helix.ChannelCustomReward
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

	twitchBinding, found := apiChannelbinding.Find(channel, platformentity.PlatformTwitch)
	if !found || twitchBinding.UserID == uuid.Nil || twitchBinding.PlatformChannelID == "" {
		return CustomRewardsResult{}, nil
	}

	twitchClient, err := c.createUserClient(ctx, twitchBinding.UserID)
	if err != nil {
		return CustomRewardsResult{}, fmt.Errorf("failed to create twitch client: %w", err)
	}

	rewards, err := twitchClient.GetCustomRewards(
		&helix.GetCustomRewardsParams{
			BroadcasterID: twitchBinding.PlatformChannelID,
		},
	)
	if err != nil {
		return CustomRewardsResult{}, fmt.Errorf("cannot get custom rewards: %w", err)
	}
	if rewards.ErrorMessage != "" {
		if rewards.StatusCode == 403 && rewards.ErrorMessage == "The broadcaster must have partner or affiliate status." {
			return CustomRewardsResult{
				Rewards:              nil,
				IsPartnerOrAffiliate: false,
			}, nil
		}
		return CustomRewardsResult{}, fmt.Errorf("cannot get custom rewards: %s", rewards.ErrorMessage)
	}

	return CustomRewardsResult{
		Rewards:              rewards.Data.ChannelCustomRewards,
		IsPartnerOrAffiliate: true,
	}, nil
}
