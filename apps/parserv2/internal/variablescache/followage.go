package variablescache

import "github.com/nicklaw5/helix"

func (c *VariablesCacheService) GetFollowAge() *helix.UserFollow {
	c.locks.twitchFollow.Lock()
	defer c.locks.twitchFollow.Unlock()

	if c.cache.TwitchFollow != nil {
		return c.cache.TwitchFollow
	}

	follow, err := c.Services.Twitch.Client.GetUsersFollows(&helix.UsersFollowsParams{
		FromID: c.Context.SenderId,
		ToID:   c.Context.ChannelId,
	})

	if err == nil && len(follow.Data.Follows) != 0 {
		c.cache.TwitchFollow = &follow.Data.Follows[0]
	}

	return c.cache.TwitchFollow
}
