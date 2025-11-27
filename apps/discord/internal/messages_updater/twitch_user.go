package messages_updater

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
)

func (c *MessagesUpdater) getTwitchUser(userId string) (helix.User, error) {
	users, err := c.twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: []string{userId},
		},
	)
	if err != nil {
		return helix.User{}, err
	}

	if users == nil || len(users.Data.Users) == 0 {
		return helix.User{}, errors.New("user not found")
	}

	twitchUser := users.Data.Users[0]
	return twitchUser, nil
}
