package variables_cache

import (
	"regexp"
	"sync"

	"github.com/satont/tsuwari/apps/parser/internal/config/twitch"
	usersauth "github.com/satont/tsuwari/apps/parser/internal/twitch/user"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/dota"
	"github.com/satont/tsuwari/libs/grpc/generated/eval"

	"github.com/go-redis/redis/v9"
	"github.com/satont/go-helix/v2"
	"gorm.io/gorm"
)

type variablesCache struct {
	Stream       *model.ChannelsStreams
	DbUserStats  *model.UsersStats
	TwitchUser   *helix.User
	TwitchFollow *helix.UserFollow
	Integrations []model.ChannelsIntegrations
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
	BotsGrpc  bots.BotsClient
	DotaGrpc  dota.DotaClient
	EvalGrpc  eval.EvalClient
}

type ExecutionContext struct {
	ChannelName string
	ChannelId   string
	SenderId    string
	SenderName  string
	Text        *string
	Services    ExecutionServices
	IsCommand   bool
	Command     *model.ChannelsCommands
}

type VariablesCacheService struct {
	ExecutionContext
	cache variablesCache
	locks *variablesLocks
}

type VariablesCacheOpts struct {
	Text        *string
	SenderId    string
	ChannelName string
	ChannelId   string
	SenderName  *string
	Redis       *redis.Client
	Regexp      *regexp.Regexp
	Twitch      *twitch.Twitch
	DB          *gorm.DB
	IsCommand   bool
	Command     *model.ChannelsCommands
	BotsGrpc    bots.BotsClient
	DotaGrpc    dota.DotaClient
	EvalGrpc    eval.EvalClient
}

func New(opts VariablesCacheOpts) *VariablesCacheService {
	cache := &VariablesCacheService{
		ExecutionContext: ExecutionContext{
			ChannelName: opts.ChannelName,
			ChannelId:   opts.ChannelId,
			SenderId:    opts.SenderId,
			SenderName:  *opts.SenderName,
			Text:        opts.Text,
			Services: ExecutionServices{
				Redis:    opts.Redis,
				Regexp:   opts.Regexp,
				Twitch:   opts.Twitch,
				Db:       opts.DB,
				BotsGrpc: opts.BotsGrpc,
				DotaGrpc: opts.DotaGrpc,
				EvalGrpc: opts.EvalGrpc,
			},
			IsCommand: opts.IsCommand,
			Command:   opts.Command,
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
