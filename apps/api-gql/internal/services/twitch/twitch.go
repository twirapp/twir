package twitch

import (
	config "github.com/satont/twir/libs/config"
	twitchcahe "github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TokensClient       tokens.TokensClient
	Config             config.Config
	CachedTwitchClient *twitchcahe.CachedTwitchClient
}

func New(opts Opts) *Service {
	return &Service{
		tokensClient:       opts.TokensClient,
		config:             opts.Config,
		cachedTwitchClient: opts.CachedTwitchClient,
	}
}

type Service struct {
	tokensClient       tokens.TokensClient
	config             config.Config
	cachedTwitchClient *twitchcahe.CachedTwitchClient
}
