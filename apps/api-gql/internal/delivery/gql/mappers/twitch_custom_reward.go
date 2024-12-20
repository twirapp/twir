package mappers

import (
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
)

func TwitchCustomRewardTo(reward helix.ChannelCustomReward) gqlmodel.TwirTwitchChannelReward {
	model := gqlmodel.TwirTwitchChannelReward{
		ID:                  reward.ID,
		BroadcasterName:     reward.BroadcasterName,
		BroadcasterLogin:    reward.BroadcasterLogin,
		BroadcasterID:       reward.BroadcasterID,
		Image:               nil,
		BackgroundColor:     reward.BackgroundColor,
		IsEnabled:           reward.IsEnabled,
		Cost:                reward.Cost,
		Title:               reward.Title,
		Prompt:              reward.Prompt,
		IsUserInputRequired: reward.IsUserInputRequired,
		MaxPerStreamSetting: &gqlmodel.TwirTwitchChannelRewardMaxPerStreamSetting{
			IsEnabled:    reward.MaxPerStreamSetting.IsEnabled,
			MaxPerStream: reward.MaxPerStreamSetting.MaxPerStream,
		},
		MaxPerUserPerStreamSetting: &gqlmodel.TwirTwitchChannelRewardMaxPerUserPerStreamSetting{
			IsEnabled:           reward.MaxPerUserPerStreamSetting.IsEnabled,
			MaxPerUserPerStream: reward.MaxPerUserPerStreamSetting.MaxPerUserPerStream,
		},
		GlobalCooldownSetting: &gqlmodel.TwirTwitchChannelRewardGlobalCooldownSetting{
			IsEnabled:             reward.GlobalCooldownSetting.IsEnabled,
			GlobalCooldownSeconds: reward.GlobalCooldownSetting.GlobalCooldownSeconds,
		},
		IsPaused:                          reward.IsPaused,
		IsInStock:                         reward.IsInStock,
		ShouldRedemptionsSkipRequestQueue: reward.ShouldRedemptionsSkipRequestQueue,
		RedemptionsRedeemedCurrentStream:  reward.RedemptionsRedeemedCurrentStream,
		CooldownExpiresAt:                 reward.CooldownExpiresAt,
	}
	var image *gqlmodel.TwirTwitchChannelRewardImage
	if reward.Image.Url1x == "" {
		image = &gqlmodel.TwirTwitchChannelRewardImage{
			URL1x: reward.DefaultImage.Url1x,
			URL2x: reward.DefaultImage.Url2x,
			URL4x: reward.DefaultImage.Url4x,
		}
	} else {
		image = &gqlmodel.TwirTwitchChannelRewardImage{
			URL1x: reward.Image.Url1x,
			URL2x: reward.Image.Url2x,
			URL4x: reward.Image.Url4x,
		}
	}
	model.Image = image

	return model
}
