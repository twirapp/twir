package mappers

import (
	"time"

	"github.com/kvizyx/twitchy/helix"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
)

func TwitchCustomRewardTo(reward helix.CustomReward) gqlmodel.TwirTwitchChannelReward {
	redemptionsRedeemedCurrentStream := 0
	if reward.RedemptionsRedeemedCurrentStream != nil {
		redemptionsRedeemedCurrentStream = *reward.RedemptionsRedeemedCurrentStream
	}
	cooldownExpiresAt := ""
	if reward.CooldownExpiresAt != nil {
		cooldownExpiresAt = reward.CooldownExpiresAt.Format(time.RFC3339Nano)
	}
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
			MaxPerStream: int(reward.MaxPerStreamSetting.MaxPerStream),
		},
		MaxPerUserPerStreamSetting: &gqlmodel.TwirTwitchChannelRewardMaxPerUserPerStreamSetting{
			IsEnabled:           reward.MaxPerUserPerStreamSetting.IsEnabled,
			MaxPerUserPerStream: int(reward.MaxPerUserPerStreamSetting.MaxPerUserPerStream),
		},
		GlobalCooldownSetting: &gqlmodel.TwirTwitchChannelRewardGlobalCooldownSetting{
			IsEnabled:             reward.GlobalCooldownSetting.IsEnabled,
			GlobalCooldownSeconds: int(reward.GlobalCooldownSetting.GlobalCooldownSeconds),
		},
		IsPaused:                          reward.IsPaused,
		IsInStock:                         reward.IsInStock,
		ShouldRedemptionsSkipRequestQueue: reward.ShouldRedemptionsSkipRequestQueue,
		RedemptionsRedeemedCurrentStream:  redemptionsRedeemedCurrentStream,
		CooldownExpiresAt:                 cooldownExpiresAt,
	}
	var image *gqlmodel.TwirTwitchChannelRewardImage
	if reward.Image == nil || reward.Image.URL1x == "" {
		image = &gqlmodel.TwirTwitchChannelRewardImage{
			URL1x: reward.DefaultImage.URL1x,
			URL2x: reward.DefaultImage.URL2x,
			URL4x: reward.DefaultImage.URL4x,
		}
	} else {
		image = &gqlmodel.TwirTwitchChannelRewardImage{
			URL1x: reward.Image.URL1x,
			URL2x: reward.Image.URL2x,
			URL4x: reward.Image.URL4x,
		}
	}
	model.Image = image

	return model
}
