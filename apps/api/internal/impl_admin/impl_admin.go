package impl_admin

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/api/internal/impl_admin/badges"
	"github.com/satont/twir/apps/api/internal/impl_admin/users"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/types/types/api/modules"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/discord"
	integrationsGrpc "github.com/twirapp/twir/libs/grpc/integrations"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Admin struct {
	*users.Users
	*badges.Badges
}

type Opts struct {
	fx.In

	Redis          *redis.Client
	DB             *gorm.DB
	Config         config.Config
	SessionManager *scs.SessionManager

	TokensGrpc        tokens.TokensClient
	IntegrationsGrpc  integrationsGrpc.IntegrationsClient
	ParserGrpc        parser.ParserClient
	WebsocketsGrpc    websockets.WebsocketClient
	DiscordGrpc       discord.DiscordClient
	Logger            logger.Logger
	Bus               *buscore.Bus
	TTSSettingsCacher *generic_cacher.GenericCacher[modules.TTSSettings]
}

func New(opts Opts) *Admin {
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
		Logger:            opts.Logger,
		Bus:               opts.Bus,
		TTSSettingsCacher: opts.TTSSettingsCacher,
	}

	return &Admin{
		Users:  &users.Users{Deps: d},
		Badges: badges.NewBadges(d),
	}
}
