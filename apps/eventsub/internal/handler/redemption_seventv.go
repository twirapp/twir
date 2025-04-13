package handler

import (
	"context"
	"errors"
	"slices"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/seventv"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (c *Handler) handleRewardsSevenTvEmote(
	event *eventsub_bindings.EventChannelPointsRewardRedemptionAdd,
) error {
	if c.config.SevenTvToken == "" || event.UserInput == "" {
		return nil
	}

	settings := &model.ChannelsIntegrationsSettingsSeventv{}
	err := c.gorm.
		Where(`"channel_id" = ?`, event.BroadcasterUserID).
		First(settings).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	if event.Reward.ID != settings.RewardIdForRemoveEmote.String &&
		event.Reward.ID != settings.RewardIdForAddEmote.String {
		return nil
	}

	ctx := context.TODO()

	client := seventv.NewClient(c.config.SevenTvToken)

	broadcasterProfile, err := client.GetProfileByTwitchId(ctx, event.BroadcasterUserID)
	if err != nil {
		return err
	}

	if broadcasterProfile.Users.UserByConnection.Style.ActiveEmoteSet == nil {
		return nil
	}

	emote, err := client.GetOneEmoteByNameOrLink(ctx, event.UserInput)
	if err != nil {
		return err
	}

	emoteId := emote.Id

	if event.Reward.ID == settings.RewardIdForRemoveEmote.String {
		if settings.DeleteEmotesOnlyAddedByApp && !slices.Contains(settings.AddedEmotes, emoteId) {
			return nil
		}

		err = client.RemoveEmote(
			ctx,
			broadcasterProfile.Users.UserByConnection.Style.ActiveEmoteSet.Id,
			event.UserInput,
			emoteId,
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

		err = c.gorm.Save(settings).Error

		return err
	}

	if event.Reward.ID == settings.RewardIdForAddEmote.String {
		err = client.AddEmote(
			ctx,
			broadcasterProfile.Users.UserByConnection.Style.ActiveEmoteSet.Id,
			emoteId,
			emote.DefaultName,
		)
		if err != nil {
			return err
		}

		settings.AddedEmotes = append(settings.AddedEmotes, emoteId)
		err = c.gorm.Save(settings).Error
		return err
	}

	return nil
}
