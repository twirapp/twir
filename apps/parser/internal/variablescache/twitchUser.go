package variables_cache

import (
	"github.com/samber/do"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/parser/internal/config/twitch"
	"github.com/satont/tsuwari/apps/parser/internal/di"
)

func (c *VariablesCacheService) GetTwitchUser() *helix.User {
	twitchClient := do.MustInvoke[twitch.Twitch](di.Provider)

	c.locks.twitchUser.Lock()
	defer c.locks.twitchUser.Unlock()

	if c.cache.TwitchUser != nil {
		return c.cache.TwitchUser
	}

	users, err := twitchClient.Client.GetUsers(&helix.UsersParams{
		IDs: []string{c.SenderId},
	})

	if err == nil && len(users.Data.Users) != 0 {
		c.cache.TwitchUser = &users.Data.Users[0]
	}

	return c.cache.TwitchUser
}
