package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) RaidChannel(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	hydratedString, hydratedErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		operation.Input.String,
		data,
	)
	if hydratedErr != nil {
		return hydratedErr
	}
	if hydratedString == "" {
		return nil
	}

	if data.ChannelID == "" {
		return errors.New("no channel id provided")
	}

	twitchClient, twitchClientErr := c.getHelixChannelApiClient(context.TODO(), data.ChannelID)
	if twitchClientErr != nil {
		return twitchClientErr
	}

	user, userErr := c.getHelixUserByLogin(twitchClient, hydratedString)
	if userErr != nil {
		return userErr
	}

	raidReq, raidErr := twitchClient.StartRaid(
		&helix.StartRaidParams{
			FromBroadcasterID: data.ChannelID,
			ToBroadcasterID:   user.ID,
		},
	)

	if raidErr != nil {
		return raidErr
	}
	if raidReq.ErrorMessage != "" {
		return fmt.Errorf("cannot start raid: %s", raidReq.ErrorMessage)
	}

	return nil
}
