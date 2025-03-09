package impl_unprotected

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/modules"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/twitch"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	apimodules "github.com/satont/twir/libs/types/types/api/modules"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/discord"
	integrationsGrpc "github.com/twirapp/twir/libs/grpc/integrations"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type UnProtected struct {
	*twitch.Twitch
	*modules.Modules
}

type Opts struct {
	fx.In

	Redis          *redis.Client
	DB             *gorm.DB
	Config         cfg.Config
	SessionManager *scs.SessionManager

	IntegrationsGrpc integrationsGrpc.IntegrationsClient
	TokensGrpc       tokens.TokensClient
	ParserGrpc       parser.ParserClient
	DiscordGrpc      discord.DiscordClient

	Bus               *buscore.Bus
	Logger            logger.Logger
	TTSSettingsCacher *generic_cacher.GenericCacher[apimodules.TTSSettings]
}

func New(opts Opts) *UnProtected {
	d := &impl_deps.Deps{
		Redis:          opts.Redis,
		Db:             opts.DB,
		Config:         opts.Config,
		SessionManager: opts.SessionManager,
		Grpc: &impl_deps.Grpc{
			Tokens:       opts.TokensGrpc,
			Integrations: opts.IntegrationsGrpc,
			Parser:       opts.ParserGrpc,
			Discord:      opts.DiscordGrpc,
		},
		Bus:               opts.Bus,
		Logger:            opts.Logger,
		TTSSettingsCacher: opts.TTSSettingsCacher,
	}

	return &UnProtected{
		Twitch: &twitch.Twitch{
			Deps: d,
		},
		Modules: &modules.Modules{
			Deps: d,
		},
	}
}
