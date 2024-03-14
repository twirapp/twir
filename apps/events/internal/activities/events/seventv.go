package events

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/seventv"
	"go.temporal.io/sdk/activity"
	"gorm.io/gorm"
)

func (c *Activity) SevenTvEmoteManage(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	hydratedString, hydrateErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		operation.Input.String,
		data,
	)
	if hydrateErr != nil {
		return fmt.Errorf("cannot hydrate string %w", hydrateErr)
	}
	if len(hydratedString) == 0 {
		return nil
	}

	broadcasterProfile, err := seventv.GetProfile(ctx, data.ChannelID)
	if err != nil {
		return err
	}

	if broadcasterProfile.EmoteSet == nil {
		return nil
	}

	emoteId := seventv.FindEmoteIdInInput(hydratedString)
	if emoteId == "" {
		return nil
	}

	settings := &model.ChannelsIntegrationsSettingsSeventv{}
	err = c.db.
		WithContext(ctx).
		Where(`"channel_id" = ?`, data.ChannelID).
		First(settings).
		Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	if operation.Type == model.OperationSevenTvAddEmote {
		err := seventv.AddEmote(
			ctx,
			c.cfg.SevenTvToken,
			hydratedString,
			broadcasterProfile.EmoteSet.Id,
		)
		if err != nil {
			return err
		}

		settings.AddedEmotes = append(settings.AddedEmotes, emoteId)
		err = c.db.WithContext(ctx).Save(settings).Error
		if err != nil {
			return err
		}

		return nil
	} else {
		if settings.DeleteEmotesOnlyAddedByApp && !slices.Contains(settings.AddedEmotes, emoteId) {
			return nil
		}

		err := seventv.RemoveEmote(
			ctx,
			c.cfg.SevenTvToken,
			hydratedString,
			broadcasterProfile.EmoteSet.Id,
		)
		if err != nil {
			return err
		}

		settings.AddedEmotes = lo.Filter(
			settings.AddedEmotes,
			func(s string, _ int) bool {
				return s != emoteId
			},
		)

		err = c.db.WithContext(ctx).Save(settings).Error
		if err != nil {
			return err
		}

		return nil
	}
}
