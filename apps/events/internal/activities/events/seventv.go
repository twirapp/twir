package events

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/events/internal/shared"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/seventv"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.temporal.io/sdk/activity"
	"gorm.io/gorm"
)

func (c *Activity) SevenTvEmoteManage(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if operation.Input == nil || *operation.Input == "" {
		return fmt.Errorf("input is required for SevenTV emote manage operation")
	}

	hydratedString, hydrateErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		*operation.Input,
		data,
	)
	if hydrateErr != nil {
		return fmt.Errorf("cannot hydrate string %w", hydrateErr)
	}
	if len(hydratedString) == 0 {
		return nil
	}

	client := seventv.NewClient(c.cfg.SevenTvToken)

	broadcasterProfile, err := client.GetProfileByTwitchId(ctx, data.ChannelID)
	if err != nil {
		return err
	}

	if broadcasterProfile.Users.UserByConnection.Style.ActiveEmoteSet == nil {
		return nil
	}

	emote, err := client.GetOneEmoteByNameOrLink(ctx, hydratedString)
	if err != nil {
		return err
	}

	settings := &deprecatedgormmodel.ChannelsIntegrationsSettingsSeventv{}
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

	if operation.Type == model.EventOperationTypeSeventvAddEmote {
		err = client.AddEmote(
			ctx,
			broadcasterProfile.Users.UserByConnection.Style.ActiveEmoteSet.Id,
			emote.Id,
			emote.DefaultName,
		)
		if err != nil {
			return err
		}

		settings.AddedEmotes = append(settings.AddedEmotes, emote.Id)
		err = c.db.WithContext(ctx).Save(settings).Error
		if err != nil {
			return err
		}

		return nil
	}

	if settings.DeleteEmotesOnlyAddedByApp && !slices.Contains(settings.AddedEmotes, emote.Id) {
		return nil
	}

	err = client.RemoveEmote(
		ctx,
		broadcasterProfile.Users.UserByConnection.Style.ActiveEmoteSet.Id,
		emote.DefaultName,
		emote.Id,
	)
	if err != nil {
		return err
	}

	settings.AddedEmotes = lo.Filter(
		settings.AddedEmotes,
		func(s string, _ int) bool {
			return s != emote.Id
		},
	)

	err = c.db.WithContext(ctx).Save(settings).Error
	if err != nil {
		return err
	}

	return nil
}
