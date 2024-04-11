package resolvers

import (
	"github.com/nicklaw5/helix/v2"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/apps/api-gql/internal/sessions"
	twitchcahe "github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	sessions           *sessions.Sessions
	gorm               *gorm.DB
	twitchClient       *helix.Client
	cachedTwitchClient *twitchcahe.CachedTwitchClient
}

type Opts struct {
	fx.In

	Sessions           *sessions.Sessions
	Gorm               *gorm.DB
	Config             config.Config
	TokensGrpc         tokens.TokensClient
	CachedTwitchClient *twitchcahe.CachedTwitchClient
}

func New(opts Opts) (*Resolver, error) {
	twitchClient, err := twitch.NewAppClient(opts.Config, opts.TokensGrpc)
	if err != nil {
		return nil, err
	}

	return &Resolver{
		sessions:           opts.Sessions,
		gorm:               opts.Gorm,
		twitchClient:       twitchClient,
		cachedTwitchClient: opts.CachedTwitchClient,
	}, nil
}
