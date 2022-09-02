package variablescache

import "github.com/nicklaw5/helix"

func (c *VariablesCacheService) GetTwitchUser() *helix.User {
	c.locks.twitchUser.Lock()
	defer c.locks.twitchUser.Unlock()

	if c.cache.TwitchUser != nil {
		return c.cache.TwitchUser
	}

	users, err := c.Services.Twitch.Client.GetUsers(&helix.UsersParams{
		IDs: []string{c.Context.SenderId},
	})

	if err == nil && len(users.Data.Users) != 0 {
		c.cache.TwitchUser = &users.Data.Users[0]
	}

	return c.cache.TwitchUser
}
