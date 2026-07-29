package cacher

import (
	"context"
	"github.com/kvizyx/twitchy/helix"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
)

// GetFollowAge implements types.VariablesCacher
func (c *cacher) GetTwitchUserFollow(ctx context.Context, userID string) *helix.ChannelFollower {
	c.locks.twitchFollow.Lock()
	defer c.locks.twitchFollow.Unlock()

	if c.cache.twitchUserFollows[userID] != nil {
		return c.cache.twitchUserFollows[userID]
	}

	dbChannel, err := c.getDbChannel(ctx)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}
	if dbChannel.BotID == "" {
		return nil
	}

	twitchClient, err := twitch.NewChannelBotClientWithContext(
		ctx,
		dbChannel.BotID,
		c.parseCtxChannel.ID,
		*c.services.Config,
	)
	if err != nil {
		return nil
	}

	follow, err := twitchClient.Channels.GetChannelFollowers(
		ctx,
		helix.GetChannelFollowersRequest{
			BroadcasterID: c.parseCtxChannel.ID,
			UserID:        &userID,
		},
	)
	if err != nil {
		return nil
	}

	if len(follow.Data) != 0 {
		c.cache.twitchUserFollows[userID] = &follow.Data[0]
	}

	return c.cache.twitchUserFollows[userID]
}

// GetGbUser implements types.VariablesCacher
func (c *cacher) GetGbUserStats(ctx context.Context, userId string) *model.UsersStats {
	c.locks.dbUserStats.Lock()
	defer c.locks.dbUserStats.Unlock()

	if c.cache.dbUserStats != nil {
		return c.cache.dbUserStats
	}

	result := &model.UsersStats{}
	dbChannel, err := c.getDbChannel(ctx)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	err = c.services.Gorm.
		WithContext(ctx).
		Where(`user_id = ? AND channel_id = ?::uuid`, userId, dbChannel.ChannelID).
		Find(result).
		Error
	if err == nil {
		c.cache.dbUserStats = result
	}

	return c.cache.dbUserStats
}

// GetTwitchChannel implements types.VariablesCacher
func (c *cacher) GetTwitchChannel(ctx context.Context) *helix.ChannelInformation {
	c.locks.twitchChannel.Lock()
	defer c.locks.twitchChannel.Unlock()

	if c.cache.twitchChannel != nil {
		return c.cache.twitchChannel
	}

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		*c.services.Config,
		c.services.Bus,
	)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	channel, err := twitchClient.Channels.GetChannelInformation(
		ctx,
		helix.GetChannelInformationRequest{
			BroadcasterIDs: []string{c.parseCtxChannel.ID},
		},
	)

	if err == nil && len(channel.Data) != 0 {
		c.cache.twitchChannel = &channel.Data[0]
	}

	return c.cache.twitchChannel
}

// GetTwitchSenderUser implements types.VariablesCacher
func (c *cacher) GetTwitchSenderUser(ctx context.Context) *helix.User {
	user, err := c.GetTwitchUserById(ctx, c.parseCtxSender.ID)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	return user
}

// GetTwitchUserById implements types.VariablesCacher
func (c *cacher) GetTwitchUserById(ctx context.Context, userId string) (*helix.User, error) {
	c.locks.cachedTwitchUsersById.Lock()
	defer c.locks.cachedTwitchUsersById.Unlock()

	if user, ok := c.cache.cachedTwitchUsersById[userId]; ok {
		return user, nil
	}

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		*c.services.Config,
		c.services.Bus,
	)
	if err != nil {
		return nil, err
	}

	users, err := twitchClient.Users.GetUsers(
		ctx,
		helix.GetUsersRequest{
			IDs: []string{userId},
		},
	)
	if err != nil {
		return nil, err
	}
	if len(users.Data) == 0 {
		return nil, nil
	}

	c.cache.cachedTwitchUsersById[userId] = &users.Data[0]
	c.cache.cachedTwitchUsersByName[users.Data[0].Login] = &users.Data[0]

	return &users.Data[0], nil
}

func (c *cacher) GetTwitchUserByName(ctx context.Context, userName string) (*helix.User, error) {
	c.locks.cachedTwitchUsersByName.Lock()
	defer c.locks.cachedTwitchUsersByName.Unlock()

	if user, ok := c.cache.cachedTwitchUsersByName[userName]; ok {
		return user, nil
	}

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		*c.services.Config,
		c.services.Bus,
	)
	if err != nil {
		return nil, err
	}

	users, err := twitchClient.Users.GetUsers(
		ctx,
		helix.GetUsersRequest{
			Logins: []string{userName},
		},
	)
	if err != nil {
		return nil, err
	}
	if len(users.Data) == 0 {
		return nil, nil
	}

	c.cache.cachedTwitchUsersByName[userName] = &users.Data[0]
	c.cache.cachedTwitchUsersById[users.Data[0].ID] = &users.Data[0]

	return &users.Data[0], nil
}
