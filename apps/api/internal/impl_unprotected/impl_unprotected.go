package impl_unprotected

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/auth"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/commands"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/modules"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/songs"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/stats"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/twitch"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/eventsub"
	integrationsGrpc "github.com/satont/twir/libs/grpc/generated/integrations"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"github.com/satont/twir/libs/grpc/generated/scheduler"
	"github.com/satont/twir/libs/grpc/generated/timers"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type UnProtected struct {
	*twitch.Twitch
	*stats.Stats
	*commands.Commands
	*auth.Auth
	*modules.Modules
	*songs.Songs
}

type Opts struct {
	fx.In

	Redis          *redis.Client
	DB             *gorm.DB
	Config         *cfg.Config
	SessionManager *scs.SessionManager

	IntegrationsGrpc integrationsGrpc.IntegrationsClient
	TokensGrpc       tokens.TokensClient
	BotsGrpc         bots.BotsClient
	ParserGrpc       parser.ParserClient
	SchedulerGrpc    scheduler.SchedulerClient
	TimersGrpc       timers.TimersClient
	EventSubGrpc     eventsub.EventSubClient

	Logger logger.Logger
}

func New(opts Opts) *UnProtected {
	d := &impl_deps.Deps{
		Redis:          opts.Redis,
		Db:             opts.DB,
		Config:         opts.Config,
		SessionManager: opts.SessionManager,
		Grpc: &impl_deps.Grpc{
			Tokens:       opts.TokensGrpc,
			Bots:         opts.BotsGrpc,
			Integrations: opts.IntegrationsGrpc,
			Parser:       opts.ParserGrpc,
			Scheduler:    opts.SchedulerGrpc,
			Timers:       opts.TimersGrpc,
			EventSub:     opts.EventSubGrpc,
		},
		Logger: opts.Logger,
	}

	return &UnProtected{
		Twitch: &twitch.Twitch{
			Deps: d,
		},
		Stats: &stats.Stats{
			Deps: d,
		},
		Commands: &commands.Commands{
			Deps: d,
		},
		Auth: &auth.Auth{
			Deps: d,
			TwitchScopes: []string{
				"moderation:read",
				"channel:manage:broadcast",
				"channel:read:redemptions",
				"moderator:read:chatters",
				"moderator:manage:shoutouts",
				"moderator:manage:banned_users",
				"channel:read:vips",
				"channel:manage:vips",
				"channel:manage:moderators",
				"moderator:read:followers",
				"moderator:manage:chat_settings",
				"channel:read:polls",
				"channel:manage:polls",
				"channel:read:predictions",
				"channel:manage:predictions",
				"channel:read:subscriptions",
				"channel:moderate",
			},
		},
		Modules: &modules.Modules{
			Deps: d,
		},
		Songs: &songs.Songs{
			Deps: d,
		},
	}
}
