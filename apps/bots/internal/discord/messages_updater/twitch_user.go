package discordmessagesupdater

import (
	"context"
	"errors"

	"github.com/kvizyx/twitchy/helix"
)

func (c *MessagesUpdater) getTwitchUser(ctx context.Context, userID string) (helix.User, error) {
	users, err := c.twitchClient.Users.GetUsers(
		ctx,
		helix.GetUsersRequest{
			IDs: []string{userID},
		},
	)
	if err != nil {
		return helix.User{}, err
	}

	if users == nil || len(users.Data) == 0 {
		return helix.User{}, errors.New("user not found")
	}

	twitchUser := users.Data[0]
	return twitchUser, nil
}
