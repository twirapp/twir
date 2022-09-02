package variablescache

import (
	"regexp"
	"sync"
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
	Integrations *[]model.ChannelInegrationWithRelation
	FaceitData   *FaceitResult
}

type variablesLocks struct {
	stream            *sync.Mutex
	dbUser            *sync.Mutex
	twitchUser        *sync.Mutex
	twitchFollow      *sync.Mutex
	integrations      *sync.Mutex
	faceitIntegration *sync.Mutex
	faceitMatches     *sync.Mutex
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
			faceitMatches:     &sync.Mutex{},
		},
	}

	return cache
}
