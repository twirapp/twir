package impl_protected

import (
	"log/slog"

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
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	commandwithgroupandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	channelseventsmodel "github.com/twirapp/twir/libs/repositories/events/model"
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
	Logger            *slog.Logger
	SpotifyRepository channelsintegrationsspotify.Repository

	Redis          *redis.Client
	DB             *gorm.DB
	SessionManager *scs.SessionManager

	Bus                               *buscore.Bus
	Config                            config.Config
	ChannelsEventsWithOperationsCache *generic_cacher.GenericCacher[[]channelseventsmodel.Event]
	CommandsCache                     *generic_cacher.GenericCacher[[]commandwithgroupandresponsesmodel.CommandWithGroupAndResponses]
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
		SpotifyRepo:                       opts.SpotifyRepository,
		ChannelsEventsWithOperationsCache: opts.ChannelsEventsWithOperationsCache,
		ChannelsCommandsCache:             opts.CommandsCache,
	}

	return &Protected{
		Integrations: &integrations.Integrations{Deps: d},
		Modules:      &modules.Modules{Deps: d},
		Twitch:       &twitch.Twitch{Deps: d},
		Overlays:     &overlays.Overlays{Deps: d},
	}
}
