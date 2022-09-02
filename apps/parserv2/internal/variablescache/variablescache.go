package variablescache

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sync"
	"time"
	"tsuwari/parser/internal/config/twitch"
	model "tsuwari/parser/internal/models"
	"tsuwari/parser/internal/variables/stream"

	"github.com/go-redis/redis/v9"
	"github.com/nicklaw5/helix"
	"gorm.io/gorm"
)

type VariablesCacheServices struct {
	Redis  *redis.Client
	Regexp regexp.Regexp
	Twitch *twitch.Twitch
	Db     *gorm.DB
}

type VariablesCacheContext struct {
	ChannelId  string
	SenderId   string
	SenderName string
	Text       string
}

type variablesCache struct {
	Stream       *stream.HelixStream
	DbUserStats  *model.UsersStats
	TwitchUser   *helix.User
	TwitchFollow *helix.UserFollow
}

type variablesLocks struct {
	stream       *sync.Mutex
	dbUser       *sync.Mutex
	twitchUser   *sync.Mutex
	twitchFollow *sync.Mutex
}

type VariablesCacheService struct {
	Context  VariablesCacheContext
	Services VariablesCacheServices
	cache    variablesCache
	locks    *variablesLocks
}

func New(text string, senderId string, channelId string, senderName *string, redis *redis.Client, r regexp.Regexp, twitch *twitch.Twitch, db *gorm.DB) *VariablesCacheService {
	cache := &VariablesCacheService{
		Context: VariablesCacheContext{
			ChannelId:  channelId,
			SenderId:   senderId,
			SenderName: *senderName,
			Text:       text,
		},
		Services: VariablesCacheServices{
			Redis:  redis,
			Regexp: r,
			Twitch: twitch,
			Db:     db,
		},
		cache: variablesCache{},
		locks: &variablesLocks{
			stream:       &sync.Mutex{},
			dbUser:       &sync.Mutex{},
			twitchUser:   &sync.Mutex{},
			twitchFollow: &sync.Mutex{},
		},
	}

	return cache
}

func (c *VariablesCacheService) GetChannelStream() *stream.HelixStream {
	c.locks.stream.Lock()
	defer c.locks.stream.Unlock()

	if c.cache.Stream != nil {
		return c.cache.Stream
	}

	rCtx := context.TODO()
	rKey := "streams:" + c.Context.ChannelId
	cachedStream, _ := c.Services.Redis.Get(rCtx, rKey).Result()

	if cachedStream != "" {
		json.Unmarshal([]byte(cachedStream), &c.cache.Stream)
		return c.cache.Stream
	}

	streams, err := c.Services.Twitch.Client.GetStreams(&helix.StreamsParams{
		UserIDs: []string{c.Context.ChannelId},
	})

	if err != nil || len(streams.Data.Streams) == 0 {
		return nil
	}

	stream := stream.HelixStream{
		Stream:   streams.Data.Streams[0],
		Messages: 0,
	}

	rData, err := json.Marshal(stream)
	if err == nil {
		c.Services.Redis.Set(rCtx, rKey, rData, time.Minute*5)
	}

	c.cache.Stream = &stream
	return c.cache.Stream
}

func (c *VariablesCacheService) GetGbUser() *model.UsersStats {
	c.locks.dbUser.Lock()
	defer c.locks.dbUser.Unlock()

	if c.cache.DbUserStats != nil {
		return c.cache.DbUserStats
	}

	result := model.UsersStats{}
	err := c.Services.Db.Where(`"userId" = ? AND "channelId" = ?`, c.Context.SenderId, c.Context.ChannelId).Find(&result).Error
	if err != nil {
		c.cache.DbUserStats = &result
	} else {
		fmt.Errorf("Cannot fetch user! %v", err)
	}

	return c.cache.DbUserStats
}

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

func (c *VariablesCacheService) GetFollowAge() *helix.UserFollow {
	c.locks.twitchFollow.Lock()
	defer c.locks.twitchUser.Unlock()

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
