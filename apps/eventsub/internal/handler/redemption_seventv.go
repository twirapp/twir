package handler

import (
	"context"
	"errors"
	"slices"

	eventsub_bindings "github.com/dnsge/twitch-eventsub-bindings"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/seventv"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (c *Handler) handleRewardsSevenTvEmote(event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd) {
	if c.services.Config.SevenTvToken == "" || event.UserInput == "" {
		return
	}

	settings := &model.ChannelsIntegrationsSettingsSeventv{}
	err := c.services.Gorm.
		Where(`"channel_id" = ?`, event.BroadcasterUserID).
		First(settings).
		Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.S().Error(err)
		}
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

	emoteId := seventv.FindEmoteIdInInput(event.UserInput)
	if emoteId == "" {
		return
	}

	if event.Reward.ID == settings.RewardIdForRemoveEmote.String {
		if settings.DeleteEmotesOnlyAddedByApp && !slices.Contains(settings.AddedEmotes, emoteId) {
			return
		}

		err = seventv.RemoveEmote(
			ctx,
			c.services.Config.SevenTvToken,
			event.UserInput,
			broadcasterProfile.EmoteSet.Id,
		)
		if err != nil {
			zap.S().Error(err)
		}

		settings.AddedEmotes = lo.Filter(
			settings.AddedEmotes,
			func(s string, _ int) bool {
				return s != emoteId
			},
		)

		err = c.services.Gorm.Save(settings).Error
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
			return
		}

		settings.AddedEmotes = append(settings.AddedEmotes, emoteId)
		err = c.services.Gorm.Save(settings).Error
		if err != nil {
			zap.S().Error(err)
		}

		return
	}
}
