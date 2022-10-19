package variables_cache

import "github.com/satont/go-helix/v2"

func (c *VariablesCacheService) GetFollowAge(userId string) *helix.UserFollow {
	c.locks.twitchFollow.Lock()
	defer c.locks.twitchFollow.Unlock()

	if c.cache.TwitchFollow != nil {
		return c.cache.TwitchFollow
	}

	follow, err := c.Services.Twitch.Client.GetUsersFollows(&helix.UsersFollowsParams{
		FromID: userId,
		ToID:   c.ChannelId,
	})

	if err == nil && len(follow.Data.Follows) != 0 {
		c.cache.TwitchFollow = &follow.Data.Follows[0]
	}

	return c.cache.TwitchFollow
}
