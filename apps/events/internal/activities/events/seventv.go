package events

import (
	"context"
	"fmt"

	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/seventv"
)

func (c *Activity) SevenTvEmoteManage(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EvenData,
) error {
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

	if operation.Type == model.OperationSevenTvAddEmote {
		return seventv.AddEmote(
			ctx,
			c.cfg.SevenTvToken,
			hydratedString,
			broadcasterProfile.EmoteSet.Id,
		)
	}

	return seventv.RemoveEmote(
		ctx,
		c.cfg.SevenTvToken,
		hydratedString,
		broadcasterProfile.EmoteSet.Id,
	)
}
