package variables_cache

import (
	"regexp"
	"sync"
	"tsuwari/parser/internal/config/twitch"
	model "tsuwari/parser/internal/models"
	usersauth "tsuwari/parser/internal/twitch/user"
	"tsuwari/parser/internal/variables/stream"

	"github.com/go-redis/redis/v9"
	"github.com/nicklaw5/helix"
	"gorm.io/gorm"
)

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

type ExecutionServices struct {
	Redis     *redis.Client
	Regexp    *regexp.Regexp
	Twitch    *twitch.Twitch
	Db        *gorm.DB
	UsersAuth *usersauth.UsersTokensService
}

type ExecutionContext struct {
	ChannelId  string
	SenderId   string
	SenderName string
	Text       *string
	Services   ExecutionServices
}

type VariablesCacheService struct {
	ExecutionContext
	cache variablesCache
	locks *variablesLocks
}

type VariablesCacheOpts struct {
	Text       *string
	SenderId   string
	ChannelId  string
	SenderName *string
	Redis      *redis.Client
	Regexp     *regexp.Regexp
	Twitch     *twitch.Twitch
	DB         *gorm.DB
}

func New(opts VariablesCacheOpts) *VariablesCacheService {
	cache := &VariablesCacheService{
		ExecutionContext: ExecutionContext{
			ChannelId:  opts.ChannelId,
			SenderId:   opts.SenderId,
			SenderName: *opts.SenderName,
			Text:       opts.Text,
			Services: ExecutionServices{
				Redis:  opts.Redis,
				Regexp: opts.Regexp,
				Twitch: opts.Twitch,
				Db:     opts.DB,
			},
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
