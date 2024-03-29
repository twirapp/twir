package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) ChangeCategory(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	hydratedCategory, err := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		operation.Input.String,
		data,
	)

	if err != nil || len(hydratedCategory) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	twitchClient, twitchClientErr := c.getHelixChannelApiClient(ctx, data.ChannelID)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	category, err := twitch.SearchCategory(ctx, hydratedCategory)
	if err != nil {
		return err
	}

	editReq, err := twitchClient.EditChannelInformation(
		&helix.EditChannelInformationParams{
			BroadcasterID: data.ChannelID,
			GameID:        category.ID,
		},
	)
	if err != nil {
		return err
	}

	if editReq.ErrorMessage != "" {
		return errors.New(editReq.ErrorMessage)
	}

	return nil
}

func (c *Activity) ChangeTitle(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	hydratedTitle, err := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		operation.Input.String,
		data,
	)

	if err != nil || len(hydratedTitle) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	twitchClient, twitchClientErr := c.getHelixChannelApiClient(ctx, data.ChannelID)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	req, err := twitchClient.EditChannelInformation(
		&helix.EditChannelInformationParams{
			BroadcasterID: data.ChannelID,
			Title:         hydratedTitle,
		},
	)
	if err != nil {
		return err
	}

	if req.ErrorMessage != "" {
		return errors.New(req.ErrorMessage)
	}

	return nil
}
