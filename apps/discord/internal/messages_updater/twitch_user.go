package messages_updater

import (
	"fmt"

	"github.com/nicklaw5/helix/v2"
)

func (c *MessagesUpdater) getTwitchUser(userId string) (helix.User, error) {
	users, err := c.twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: []string{userId},
		},
	)

	if users == nil || len(users.Data.Users) == 0 {
		return helix.User{}, fmt.Errorf("user not found")
	}

	return users.Data.Users[0], err
}
