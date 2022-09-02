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
	"tsuwari/parser/pkg/helpers"

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

type VariablesCacheService struct {
	Context  VariablesCacheContext
	Services VariablesCacheServices
	Cache    VariablesCache
}

type VariablesCache struct {
	Stream       *stream.HelixStream
	DbUserStats  *model.UsersStats
	TwitchUser   *helix.User
	TwitchFollow *helix.UserFollow
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
		Cache: VariablesCache{
			Stream: nil,
		},
	}

	cache.fillCache()

	return cache
}

type MapItem struct {
	Instance string
	Func     func(wg *sync.WaitGroup)
}

func (c *VariablesCacheService) fillCache() {
	matches := c.Services.Regexp.FindAllStringSubmatch(c.Context.Text, len(c.Context.Text))
	myMap := map[string]MapItem{
		"stream": {
			Instance: "twitchStream",
			Func:     c.setChannelStream,
		},
		"user.messages": {
			Instance: "dbUser",
			Func:     c.setUser,
		},
		"user.followage": {
			Instance: "followage",
			Func:     c.setFollowAge,
		},
		"user.age": {
			Instance: "twitchUser",
			Func:     c.setTwitchUser,
		},
	}

	requesting := []string{}
	wg := sync.WaitGroup{}

	for _, match := range matches {
		if match[1] == "" {
			continue
		}

		if val, ok := myMap[match[1]]; ok {
			if helpers.Contains(requesting, val.Instance) {
				continue
			}

			requesting = append(requesting, val.Instance)
			wg.Add(1)

			go val.Func(&wg)
		}
	}

	wg.Wait()
}

func (c *VariablesCacheService) setChannelStream(wg *sync.WaitGroup) {
	defer wg.Done()

	rCtx := context.TODO()
	rKey := "streams:" + c.Context.ChannelId
	cachedStream, _ := c.Services.Redis.Get(rCtx, rKey).Result()

	if cachedStream != "" {
		json.Unmarshal([]byte(cachedStream), &c.Cache.Stream)
		return
	}

	streams, err := c.Services.Twitch.Client.GetStreams(&helix.StreamsParams{
		UserIDs: []string{c.Context.ChannelId},
	})
	if err != nil || len(streams.Data.Streams) == 0 {
		return
	}

	stream := stream.HelixStream{
		Stream:   streams.Data.Streams[0],
		Messages: 0,
	}

	rData, err := json.Marshal(stream)
	if err == nil {
		c.Services.Redis.Set(rCtx, rKey, rData, time.Minute*5)
	}

	c.Cache.Stream = &stream
}

func (c *VariablesCacheService) setUser(wg *sync.WaitGroup) {
	defer wg.Done()

	result := model.UsersStats{}
	err := c.Services.Db.Where(`"userId" = ? AND "channelId" = ?`, c.Context.SenderId, c.Context.ChannelId).Find(&result).Error
	if err != nil {
		c.Cache.DbUserStats = &result
	} else {
		fmt.Errorf("Cannot fetch user! %v", err)
	}
}

func (c *VariablesCacheService) setTwitchUser(wg *sync.WaitGroup) {
	defer wg.Done()

	users, err := c.Services.Twitch.Client.GetUsers(&helix.UsersParams{
		IDs: []string{c.Context.SenderId},
	})

	if err == nil && len(users.Data.Users) != 0 {
		c.Cache.TwitchUser = &users.Data.Users[0]
	}
}

func (c *VariablesCacheService) setFollowAge(wg *sync.WaitGroup) {
	defer wg.Done()

	follow, err := c.Services.Twitch.Client.GetUsersFollows(&helix.UsersFollowsParams{
		FromID: c.Context.SenderId,
		ToID:   c.Context.ChannelId,
	})

	if err == nil && len(follow.Data.Follows) != 0 {
		c.Cache.TwitchFollow = &follow.Data.Follows[0]
	}
}
