package impl_protected

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/api/internal/impl_deps"
	"github.com/twirapp/twir/apps/api/internal/impl_protected/integrations"
	"github.com/twirapp/twir/apps/api/internal/impl_protected/modules"
	"github.com/twirapp/twir/apps/api/internal/impl_protected/overlays"
	"github.com/twirapp/twir/apps/api/internal/impl_protected/twitch"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/discord"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/logger"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	channelseventsmodel "github.com/twirapp/twir/libs/repositories/events/model"
	apimodules "github.com/twirapp/twir/libs/types/types/api/modules"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Protected struct {
	*integrations.Integrations
	*modules.Modules
	*twitch.Twitch
	*overlays.Overlays
}

type Opts struct {
	fx.In

	WebsocketsGrpc    websockets.WebsocketClient
	DiscordGrpc       discord.DiscordClient
	Logger            logger.Logger
	SpotifyRepository channelsintegrationsspotify.Repository

	Redis          *redis.Client
	DB             *gorm.DB
	SessionManager *scs.SessionManager

	Bus                               *buscore.Bus
	TTSSettingsCacher                 *generic_cacher.GenericCacher[apimodules.TTSSettings]
	Config                            config.Config
	ChannelsEventsWithOperationsCache *generic_cacher.GenericCacher[[]channelseventsmodel.Event]
}

func New(opts Opts) *Protected {
	d := &impl_deps.Deps{
		Redis:          opts.Redis,
		Db:             opts.DB,
		Config:         opts.Config,
		SessionManager: opts.SessionManager,
		Grpc: &impl_deps.Grpc{
			Websockets: opts.WebsocketsGrpc,
			Discord:    opts.DiscordGrpc,
		},
		Logger:                            opts.Logger,
		Bus:                               opts.Bus,
		TTSSettingsCacher:                 opts.TTSSettingsCacher,
		SpotifyRepo:                       opts.SpotifyRepository,
		ChannelsEventsWithOperationsCache: opts.ChannelsEventsWithOperationsCache,
	}

	return &Protected{
		Integrations: &integrations.Integrations{Deps: d},
		Modules:      &modules.Modules{Deps: d},
		Twitch:       &twitch.Twitch{Deps: d},
		Overlays:     &overlays.Overlays{Deps: d},
	}
}
