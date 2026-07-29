package events

import (
	"context"
	"fmt"
	"strings"

	"github.com/kvizyx/twitchy/helix"
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
		data.ChannelTwitchUserID,
		data.ChannelDBID,
		*operation.Input,
		data,
	)
	if hydrateErr != nil || len(shoutoutTarget) == 0 {
		return fmt.Errorf("cannot hydrate string %w", hydrateErr)
	}

	shoutoutTarget = strings.TrimSpace(strings.ReplaceAll(shoutoutTarget, "@", ""))

	twitchClient, err := c.getHelixChannelApiClient(ctx, data.ChannelTwitchUserID)
	if err != nil {
		return err
	}

	usersReq, err := twitchClient.Users.GetUsers(
		ctx, helix.GetUsersRequest{
			Logins: []string{shoutoutTarget},
		},
	)
	if err != nil {
		return err
	}
	if len(usersReq.Data) == 0 {
		return fmt.Errorf("cannot find user with this name")
	}

	user := usersReq.Data[0]

	_, err = twitchClient.Chat.SendShoutout(
		ctx, helix.SendShoutoutRequest{
			FromBroadcasterID: twitchBroadcasterID(data),
			ToBroadcasterID:   user.ID,
			ModeratorID:       twitchBroadcasterID(data),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
