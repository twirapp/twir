package impl_unprotected

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/api/internal/impl_deps"
	"github.com/twirapp/twir/apps/api/internal/impl_unprotected/twitch"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/discord"
	"github.com/twirapp/twir/libs/logger"
	apimodules "github.com/twirapp/twir/libs/types/types/api/modules"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type UnProtected struct {
	*twitch.Twitch
}

type Opts struct {
	fx.In

	Redis          *redis.Client
	DB             *gorm.DB
	Config         cfg.Config
	SessionManager *scs.SessionManager

	DiscordGrpc discord.DiscordClient

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
			Discord: opts.DiscordGrpc,
		},
		Bus:               opts.Bus,
		Logger:            opts.Logger,
		TTSSettingsCacher: opts.TTSSettingsCacher,
	}

	return &UnProtected{
		Twitch: &twitch.Twitch{
			Deps: d,
		},
	}
}
