package twitch

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/libs/twitch"
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

	channel, err := c.channelsRepository.GetByID(ctx, parsedID)
	if err != nil {
		return CustomRewardsResult{}, fmt.Errorf("get channel: %w", err)
	}
	if channel.IsNil() || !channel.TwitchConnected() {
		return CustomRewardsResult{}, nil
	}

	twitchClient, err := twitch.NewUserClientWithContext(ctx, *channel.TwitchPlatformID, c.config, c.twirBus)
	if err != nil {
		return CustomRewardsResult{}, fmt.Errorf("failed to create twitch client: %w", err)
	}

	rewards, err := twitchClient.GetCustomRewards(
		&helix.GetCustomRewardsParams{
			BroadcasterID: *channel.TwitchPlatformID,
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
