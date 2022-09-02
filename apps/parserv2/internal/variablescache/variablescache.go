package variablescache

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"sync"
	"tsuwari/parser/internal/config/twitch"
	model "tsuwari/parser/internal/models"
	"tsuwari/parser/internal/variables/stream"

	"github.com/go-redis/redis/v9"
	"github.com/nicklaw5/helix"
	"github.com/samber/lo"
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
	Integrations *[]model.ChannelInegrationWithRelation
	FaceitData   *FaceitGame
}

type variablesLocks struct {
	stream            *sync.Mutex
	dbUser            *sync.Mutex
	twitchUser        *sync.Mutex
	twitchFollow      *sync.Mutex
	integrations      *sync.Mutex
	faceitIntegration *sync.Mutex
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
			stream:            &sync.Mutex{},
			dbUser:            &sync.Mutex{},
			twitchUser:        &sync.Mutex{},
			twitchFollow:      &sync.Mutex{},
			integrations:      &sync.Mutex{},
			faceitIntegration: &sync.Mutex{},
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
		c.Services.Redis.Set(rCtx, rKey, rData, 0)
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
	if err == nil {
		c.cache.DbUserStats = &result
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

func (c *VariablesCacheService) GetEnabledIntegrations() *[]model.ChannelInegrationWithRelation {
	c.locks.integrations.Lock()
	defer c.locks.integrations.Unlock()

	if c.cache.Integrations != nil {
		return c.cache.Integrations
	}

	result := &[]model.ChannelInegrationWithRelation{}
	err := c.Services.Db.Where(`"channelId" = ? AND enabled = ?`, c.Context.ChannelId, true).Joins("Integration").Find(result).Error

	if err == nil {
		c.cache.Integrations = result
	}

	return c.cache.Integrations
}

type FaceitGame struct {
	Lvl int `json:"skill_level"`
	Elo int `json:"faceit_elo"`
}

type FaceitResponse struct {
	Games map[string]*FaceitGame `json:"games"`
}

type FaceitDbData struct {
	Game     *string `json:"game"`
	Username string  `json:"username"`
}

func (c *VariablesCacheService) GetFaceitData() (*FaceitGame, error) {
	c.locks.faceitIntegration.Lock()
	defer c.locks.faceitIntegration.Unlock()

	if c.cache.FaceitData != nil {
		return c.cache.FaceitData, nil
	}

	integrations := c.GetEnabledIntegrations()

	if integrations == nil {
		return nil, errors.New("integrations not enabled")
	}

	integration, ok := lo.Find(*integrations, func(i model.ChannelInegrationWithRelation) bool {
		return i.Integration.Service == "FACEIT"
	})

	if !ok {
		return nil, errors.New("faceit integration not enabled")
	}

	dbData := &FaceitDbData{}
	err := json.Unmarshal([]byte(integration.Data.String), &dbData)

	if err != nil {
		return nil, errors.New("failed to read your faceit config. Are you sure you are using integration right?")
	}

	var game string

	if dbData.Game == nil {
		game = "csgo"
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://open.faceit.com/data/v4/players?nickname="+"Satonteu", nil)
	req.Header.Set("Authorization", "Bearer "+integration.Integration.APIKey.String)
	res, err := client.Do(req)

	if err != nil {
		return nil, errors.New("failed to fetch data from faceit")
	}

	data := FaceitResponse{}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, errors.New("failed to fetch data from faceit")
	}

	if data.Games[game] == nil {
		return nil, errors.New("Game " + game + " not found in faceit response.")
	}

	c.cache.FaceitData = data.Games[game]

	return data.Games[game], nil
}
