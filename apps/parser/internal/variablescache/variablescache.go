package variables_cache

import (
	"sync"

	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
)

type variablesCache struct {
	Stream        *model.ChannelsStreams
	DbUserStats   *model.UsersStats
	TwitchUser    *helix.User
	TwitchFollow  *helix.UserFollow
	TwitchChannel *helix.ChannelInformation
	Integrations  []model.ChannelsIntegrations
	FaceitData    *FaceitResult
}

type variablesLocks struct {
	stream            *sync.Mutex
	dbUser            *sync.Mutex
	twitchUser        *sync.Mutex
	twitchFollow      *sync.Mutex
	twitchChannel     *sync.Mutex
	integrations      *sync.Mutex
	faceitIntegration *sync.Mutex
	faceitMatches     *sync.Mutex
}

type ExecutionContext struct {
	ChannelName       string
	ChannelId         string
	SenderId          string
	SenderName        string
	SenderDisplayName string
	SenderBadges      []string
	Text              *string
	IsCommand         bool
	Command           *model.ChannelsCommands
}

type VariablesCacheService struct {
	ExecutionContext
	cache variablesCache
	locks *variablesLocks
}

type VariablesCacheOpts struct {
	Text              *string
	SenderId          string
	ChannelName       string
	ChannelId         string
	SenderName        *string
	SenderDisplayName *string
	DB                *gorm.DB
	IsCommand         bool
	Command           *model.ChannelsCommands
	SenderBadges      []string
}

func New(opts VariablesCacheOpts) *VariablesCacheService {
	cache := &VariablesCacheService{
		ExecutionContext: ExecutionContext{
			ChannelName:  opts.ChannelName,
			ChannelId:    opts.ChannelId,
			SenderId:     opts.SenderId,
			SenderName:   *opts.SenderName,
			SenderBadges: opts.SenderBadges,
			Text:         opts.Text,
			IsCommand:    opts.IsCommand,
			Command:      opts.Command,
		},
		cache: variablesCache{},
		locks: &variablesLocks{
			stream:            &sync.Mutex{},
			dbUser:            &sync.Mutex{},
			twitchUser:        &sync.Mutex{},
			twitchFollow:      &sync.Mutex{},
			twitchChannel:     &sync.Mutex{},
			integrations:      &sync.Mutex{},
			faceitIntegration: &sync.Mutex{},
			faceitMatches:     &sync.Mutex{},
		},
	}

	return cache
}
