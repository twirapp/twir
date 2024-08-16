package cacher

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
)

// GetChannelStream implements types.VariablesCacher
func (c *cacher) GetChannelStream(ctx context.Context) *model.ChannelsStreams {
	c.locks.stream.Lock()
	defer c.locks.stream.Unlock()

	if c.cache.stream != nil {
		return c.cache.stream
	}

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		*c.services.Config,
		c.services.GrpcClients.Tokens,
	)
	if err != nil {
		return nil
	}

	stream := model.ChannelsStreams{}

	err = c.services.Gorm.
		WithContext(ctx).
		Where(`"userId" = ?`, c.parseCtxChannel.ID).
		First(&stream).Error

	if err != nil {
		c.services.Logger.Sugar().Error(err)

		streams, err := twitchClient.GetStreams(
			&helix.StreamsParams{
				UserIDs: []string{c.parseCtxChannel.ID},
			},
		)

		if err != nil || len(streams.Data.Streams) == 0 {
			return nil
		}

		helixStream := streams.Data.Streams[0]

		tags := pq.StringArray{}
		for _, t := range helixStream.TagIDs {
			tags = append(tags, t)
		}
		stream = model.ChannelsStreams{
			ID:           helixStream.ID,
			UserId:       helixStream.UserID,
			UserLogin:    helixStream.UserLogin,
			UserName:     helixStream.UserName,
			GameId:       helixStream.GameID,
			GameName:     helixStream.GameName,
			CommunityIds: []string{},
			Type:         helixStream.Type,
			Title:        helixStream.Title,
			ViewerCount:  helixStream.ViewerCount,
			StartedAt:    helixStream.StartedAt,
			Language:     helixStream.Language,
			ThumbnailUrl: helixStream.ThumbnailURL,
			TagIds:       &tags,
			IsMature:     helixStream.IsMature,
		}

		c.services.Gorm.Save(&stream)
		c.cache.stream = &stream
	} else {
		c.cache.stream = &stream
	}

	return c.cache.stream
}

// GetFollowAge implements types.VariablesCacher
func (c *cacher) GetTwitchUserFollow(ctx context.Context, userID string) *helix.ChannelFollow {
	c.locks.twitchFollow.Lock()
	defer c.locks.twitchFollow.Unlock()

	if c.cache.twitchUserFollows[userID] != nil {
		return c.cache.twitchUserFollows[userID]
	}

	channel := model.Channels{}
	err := c.services.Gorm.
		WithContext(ctx).
		Where(`"id" = ?`, c.parseCtxChannel.ID).
		First(&channel).
		Error
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	twitchClient, err := twitch.NewBotClientWithContext(
		ctx,
		channel.BotID,
		*c.services.Config,
		c.services.GrpcClients.Tokens,
	)
	if err != nil {
		return nil
	}

	follow, err := twitchClient.GetChannelFollows(
		&helix.GetChannelFollowsParams{
			BroadcasterID: c.parseCtxChannel.ID,
			UserID:        userID,
			First:         0,
			After:         "",
		},
	)

	if follow.ErrorMessage != "" {
		fmt.Println(follow.ErrorMessage)
		return nil
	}

	if err == nil && len(follow.Data.Channels) != 0 {
		c.cache.twitchUserFollows[userID] = &follow.Data.Channels[0]
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

	err := c.services.Gorm.
		WithContext(ctx).
		Where(`"userId" = ? AND "channelId" = ?`, userId, c.parseCtxChannel.ID).
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
		c.services.GrpcClients.Tokens,
	)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return nil
	}

	channel, err := twitchClient.GetChannelInformation(
		&helix.GetChannelInformationParams{
			BroadcasterIDs: []string{c.parseCtxChannel.ID},
		},
	)

	if err == nil && len(channel.Data.Channels) != 0 {
		c.cache.twitchChannel = &channel.Data.Channels[0]
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
		c.services.GrpcClients.Tokens,
	)
	if err != nil {
		return nil, err
	}

	users, err := twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: []string{userId},
		},
	)
	if err != nil {
		return nil, err
	}
	if users.ErrorMessage != "" {
		return nil, fmt.Errorf(users.ErrorMessage)
	}

	if len(users.Data.Users) == 0 {
		return nil, nil
	}

	c.cache.cachedTwitchUsersById[userId] = &users.Data.Users[0]
	c.cache.cachedTwitchUsersByName[users.Data.Users[0].Login] = &users.Data.Users[0]

	return &users.Data.Users[0], nil
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
		c.services.GrpcClients.Tokens,
	)
	if err != nil {
		return nil, err
	}

	users, err := twitchClient.GetUsers(
		&helix.UsersParams{
			Logins: []string{userName},
		},
	)
	if err != nil {
		return nil, err
	}
	if users.ErrorMessage != "" {
		return nil, fmt.Errorf(users.ErrorMessage)
	}

	if len(users.Data.Users) == 0 {
		return nil, nil
	}

	c.cache.cachedTwitchUsersByName[userName] = &users.Data.Users[0]
	c.cache.cachedTwitchUsersById[users.Data.Users[0].ID] = &users.Data.Users[0]

	return &users.Data.Users[0], nil
}
