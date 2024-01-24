package handler

import (
	"context"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/seventv"
	"go.uber.org/zap"
)

func (c *Handler) handleRewardsSevenTvEmote(event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd) {
	if c.services.Config.SevenTvToken == "" || event.UserInput == "" {
		return
	}

	settings := &model.ChannelsIntegrationsSettingsSeventv{}
	err := c.services.Gorm.
		Where(`"channelId" = ?`, event.BroadcasterUserID).
		Find(settings).
		Error
	if err != nil {
		zap.S().Error(err)
		return
	}

	// not found
	if settings.ID.String() == "" {
		return
	}

	ctx := context.TODO()
	broadcasterProfile, err := seventv.GetProfile(ctx, event.BroadcasterUserID)
	if err != nil {
		zap.S().Error(err)
		return
	}

	if broadcasterProfile.EmoteSet == nil {
		return
	}

	if event.Reward.ID == settings.RewardIdForRemoveEmote.String {
		err = seventv.RemoveEmote(
			ctx,
			c.services.Config.SevenTvToken,
			event.UserInput,
			broadcasterProfile.EmoteSet.Id,
		)
		if err != nil {
			zap.S().Error(err)
		}

		return
	}

	if event.Reward.ID == settings.RewardIdForAddEmote.String {
		err = seventv.AddEmote(
			ctx,
			c.services.Config.SevenTvToken,
			event.UserInput,
			broadcasterProfile.EmoteSet.Id,
		)
		if err != nil {
			zap.S().Error(err)
		}

		return
	}
}
