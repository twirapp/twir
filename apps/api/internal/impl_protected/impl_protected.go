package impl_protected

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	"github.com/satont/twir/apps/api/internal/impl_protected/events"
	"github.com/satont/twir/apps/api/internal/impl_protected/files"
	"github.com/satont/twir/apps/api/internal/impl_protected/integrations"
	"github.com/satont/twir/apps/api/internal/impl_protected/modules"
	"github.com/satont/twir/apps/api/internal/impl_protected/overlays"
	"github.com/satont/twir/apps/api/internal/impl_protected/twitch"
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	apimodules "github.com/satont/twir/libs/types/types/api/modules"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/discord"
	integrationsGrpc "github.com/twirapp/twir/libs/grpc/integrations"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Protected struct {
	*integrations.Integrations
	*modules.Modules
	*events.Events
	*twitch.Twitch
	*files.Files
	*overlays.Overlays
}

type Opts struct {
	fx.In

	TokensGrpc        tokens.TokensClient
	IntegrationsGrpc  integrationsGrpc.IntegrationsClient
	ParserGrpc        parser.ParserClient
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
	ChannelsEventsWithOperationsCache *generic_cacher.GenericCacher[[]model.Event]
}

func New(opts Opts) *Protected {
	d := &impl_deps.Deps{
		Redis:          opts.Redis,
		Db:             opts.DB,
		Config:         opts.Config,
		SessionManager: opts.SessionManager,
		Grpc: &impl_deps.Grpc{
			Tokens:       opts.TokensGrpc,
			Integrations: opts.IntegrationsGrpc,
			Parser:       opts.ParserGrpc,
			Websockets:   opts.WebsocketsGrpc,
			Discord:      opts.DiscordGrpc,
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
		Events:       &events.Events{Deps: d},
		Twitch:       &twitch.Twitch{Deps: d},
		Files:        files.New(d),
		Overlays:     &overlays.Overlays{Deps: d},
	}
}
