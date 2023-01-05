package variables_cache

import (
	"github.com/samber/do"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/parser/internal/config/twitch"
	"github.com/satont/tsuwari/apps/parser/internal/di"
)

func (c *VariablesCacheService) GetFollowAge(userId string) *helix.UserFollow {
	twitchClient := do.MustInvoke[twitch.Twitch](di.Provider)

	c.locks.twitchFollow.Lock()
	defer c.locks.twitchFollow.Unlock()

	if c.cache.TwitchFollow != nil {
		return c.cache.TwitchFollow
	}

	follow, err := twitchClient.Client.GetUsersFollows(&helix.UsersFollowsParams{
		FromID: userId,
		ToID:   c.ChannelId,
	})

	if err == nil && len(follow.Data.Follows) != 0 {
		c.cache.TwitchFollow = &follow.Data.Follows[0]
	}

	return c.cache.TwitchFollow
}
