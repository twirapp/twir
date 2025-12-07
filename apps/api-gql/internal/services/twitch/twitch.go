package twitch

import (
	buscore "github.com/twirapp/twir/libs/bus-core"
	twitchcahe "github.com/twirapp/twir/libs/cache/twitch"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/repositories/users"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TwirBus            *buscore.Bus
	Config             config.Config
	CachedTwitchClient *twitchcahe.CachedTwitchClient
	UsersRepository    users.Repository
}

func New(opts Opts) *Service {
	return &Service{
		twirBus:            opts.TwirBus,
		config:             opts.Config,
		cachedTwitchClient: opts.CachedTwitchClient,
		usersRepository:    opts.UsersRepository,
	}
}

type Service struct {
	twirBus            *buscore.Bus
	config             config.Config
	cachedTwitchClient *twitchcahe.CachedTwitchClient
	usersRepository    users.Repository
}
