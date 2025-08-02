package events

import (
	"context"
	"fmt"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) ShoutoutChannel(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if operation.Input == nil || *operation.Input == "" {
		return fmt.Errorf("input is required for shoutout operation")
	}

	shoutoutTarget, hydrateErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		*operation.Input,
		data,
	)
	if hydrateErr != nil || len(shoutoutTarget) == 0 {
		return fmt.Errorf("cannot hydrate string %w", hydrateErr)
	}

	shoutoutTarget = strings.TrimSpace(strings.ReplaceAll(shoutoutTarget, "@", ""))

	twitchClient, err := c.getHelixChannelApiClient(ctx, data.ChannelID)
	if err != nil {
		return err
	}

	usersReq, err := twitchClient.GetUsers(
		&helix.UsersParams{
			Logins: []string{shoutoutTarget},
		},
	)
	if err != nil {
		return err
	}
	if usersReq.ErrorMessage != "" {
		return fmt.Errorf("cannot get user: %s", usersReq.ErrorMessage)
	}
	if len(usersReq.Data.Users) == 0 {
		return fmt.Errorf("cannot find user with this name")
	}

	user := usersReq.Data.Users[0]

	result, err := twitchClient.SendShoutout(
		&helix.SendShoutoutParams{
			FromBroadcasterID: data.ChannelID,
			ToBroadcasterID:   user.ID,
			ModeratorID:       data.ChannelID,
		},
	)
	if err != nil {
		return err
	}
	if result.ErrorMessage != "" {
		return fmt.Errorf("cannot send shoutout: %s", result.ErrorMessage)
	}

	return nil
}
