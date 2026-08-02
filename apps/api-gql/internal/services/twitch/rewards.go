package twitch

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type CustomRewardsResult struct {
	Rewards              []helix.ChannelCustomReward
	IsPartnerOrAffiliate bool
}

type ManageRewardInput struct {
	Action                            string
	ID                                string
	Title                             string
	Prompt                            string
	Cost                              int
	Enabled                           bool
	BackgroundColor                   string
	UserInputRequired                 bool
	MaxPerStreamEnabled               bool
	MaxPerStream                      int
	MaxPerUserPerStreamEnabled        bool
	MaxPerUserPerStream               int
	GlobalCooldownEnabled             bool
	GlobalCooldownSeconds             int
	ShouldRedemptionsSkipRequestQueue bool
}

func (c *Service) ManageReward(ctx context.Context, channelID string, input ManageRewardInput) (any, error) {
	parsedID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, fmt.Errorf("invalid channel id: %w", err)
	}
	channel, err := c.channelService.GetChannelByID(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("get channel: %w", err)
	}
	binding, found := channel.Binding(platformentity.PlatformTwitch)
	if !found {
		return nil, fmt.Errorf("channel has no Twitch binding")
	}
	client, err := c.createUserClient(ctx, binding.UserID)
	if err != nil {
		return nil, fmt.Errorf("create Twitch client: %w", err)
	}

	switch input.Action {
	case "create":
		response, err := client.CreateCustomReward(&helix.ChannelCustomRewardsParams{
			BroadcasterID: binding.PlatformChannelID, Title: input.Title, Prompt: input.Prompt,
			Cost: input.Cost, IsEnabled: input.Enabled, BackgroundColor: input.BackgroundColor,
			IsUserInputRequired: input.UserInputRequired, IsMaxPerStreamEnabled: input.MaxPerStreamEnabled,
			MaxPerStream: input.MaxPerStream, IsMaxPerUserPerStreamEnabled: input.MaxPerUserPerStreamEnabled,
			MaxPerUserPerStream: input.MaxPerUserPerStream, IsGlobalCooldownEnabled: input.GlobalCooldownEnabled,
			GlobalCooldownSeconds:             input.GlobalCooldownSeconds,
			ShouldRedemptionsSkipRequestQueue: input.ShouldRedemptionsSkipRequestQueue,
		})
		if err != nil {
			return nil, err
		}
		if response.ErrorMessage != "" {
			return nil, fmt.Errorf("Twitch: %s", response.ErrorMessage)
		}
		return response.Data.ChannelCustomRewards, nil
	case "update":
		response, err := client.UpdateCustomReward(&helix.UpdateChannelCustomRewardsParams{
			ID: input.ID, BroadcasterID: binding.PlatformChannelID, Title: input.Title,
			Prompt: input.Prompt, Cost: input.Cost, IsEnabled: input.Enabled,
			BackgroundColor: input.BackgroundColor, IsUserInputRequired: input.UserInputRequired,
			IsMaxPerStreamEnabled: input.MaxPerStreamEnabled, MaxPerStream: input.MaxPerStream,
			IsMaxPerUserPerStreamEnabled:      input.MaxPerUserPerStreamEnabled,
			MaxPerUserPerStream:               input.MaxPerUserPerStream,
			IsGlobalCooldownEnabled:           input.GlobalCooldownEnabled,
			GlobalCooldownSeconds:             input.GlobalCooldownSeconds,
			ShouldRedemptionsSkipRequestQueue: input.ShouldRedemptionsSkipRequestQueue,
		})
		if err != nil {
			return nil, err
		}
		if response.ErrorMessage != "" {
			return nil, fmt.Errorf("Twitch: %s", response.ErrorMessage)
		}
		return response.Data.ChannelCustomRewards, nil
	case "delete":
		response, err := client.DeleteCustomRewards(&helix.DeleteCustomRewardsParams{BroadcasterID: binding.PlatformChannelID, ID: input.ID})
		if err != nil {
			return nil, err
		}
		if response.ErrorMessage != "" {
			return nil, fmt.Errorf("Twitch: %s", response.ErrorMessage)
		}
		return map[string]bool{"deleted": true}, nil
	default:
		return nil, fmt.Errorf("unsupported action %q", input.Action)
	}
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
