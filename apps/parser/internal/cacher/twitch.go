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
			ID:             helixStream.ID,
			UserId:         helixStream.UserID,
			UserLogin:      helixStream.UserLogin,
			UserName:       helixStream.UserName,
			GameId:         helixStream.GameID,
			GameName:       helixStream.GameName,
			CommunityIds:   []string{},
			Type:           helixStream.Type,
			Title:          helixStream.Title,
			ViewerCount:    helixStream.ViewerCount,
			StartedAt:      helixStream.StartedAt,
			Language:       helixStream.Language,
			ThumbnailUrl:   helixStream.ThumbnailURL,
			TagIds:         &tags,
			IsMature:       helixStream.IsMature,
			ParsedMessages: 0,
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
func (c *cacher) GetGbUserStats(ctx context.Context) *model.UsersStats {
	c.locks.dbUserStats.Lock()
	defer c.locks.dbUserStats.Unlock()

	if c.cache.dbUserStats != nil {
		return c.cache.dbUserStats
	}

	result := &model.UsersStats{}

	err := c.services.Gorm.
		WithContext(ctx).
		Where(`"userId" = ? AND "channelId" = ?`, c.parseCtxSender.ID, c.parseCtxChannel.ID).
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

// GetTwitchUser implements types.VariablesCacher
func (c *cacher) GetTwitchSenderUser(ctx context.Context) *helix.User {
	c.locks.twitchSenderUser.Lock()
	defer c.locks.twitchSenderUser.Unlock()

	if c.cache.twitchSenderUser != nil {
		return c.cache.twitchSenderUser
	}

	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		*c.services.Config,
		c.services.GrpcClients.Tokens,
	)
	if err != nil {
		return nil
	}

	users, err := twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: []string{c.parseCtxSender.ID},
		},
	)

	if err != nil || users.ErrorMessage != "" {
		c.services.Logger.Sugar().Error(users.ErrorMessage, err)
		return nil
	}

	if err == nil && len(users.Data.Users) != 0 {
		c.cache.twitchSenderUser = &users.Data.Users[0]
	}

	return c.cache.twitchSenderUser
}
