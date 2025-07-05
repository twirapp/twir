package twitch

import (
	config "github.com/satont/twir/libs/config"
	buscore "github.com/twirapp/twir/libs/bus-core"
	twitchcahe "github.com/twirapp/twir/libs/cache/twitch"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TwirBus            *buscore.Bus
	Config             config.Config
	CachedTwitchClient *twitchcahe.CachedTwitchClient
}

func New(opts Opts) *Service {
	return &Service{
		twirBus:            opts.TwirBus,
		config:             opts.Config,
		cachedTwitchClient: opts.CachedTwitchClient,
	}
}

type Service struct {
	twirBus            *buscore.Bus
	config             config.Config
	cachedTwitchClient *twitchcahe.CachedTwitchClient
}
