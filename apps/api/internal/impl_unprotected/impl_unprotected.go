package impl_unprotected

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/auth"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/commands"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/community"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/modules"
	public_settings "github.com/satont/twir/apps/api/internal/impl_unprotected/public-settings"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/songs"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/stats"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/tts"
	"github.com/satont/twir/apps/api/internal/impl_unprotected/twitch"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/bots"
	"github.com/twirapp/twir/libs/grpc/discord"
	"github.com/twirapp/twir/libs/grpc/eventsub"
	integrationsGrpc "github.com/twirapp/twir/libs/grpc/integrations"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/scheduler"
	"github.com/twirapp/twir/libs/grpc/timers"
	"github.com/twirapp/twir/libs/grpc/tokens"
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
	*tts.Tts
	*community.Community
	*public_settings.PublicSettings
}

type Opts struct {
	fx.In

	Redis          *redis.Client
	DB             *gorm.DB
	Config         cfg.Config
	SessionManager *scs.SessionManager

	IntegrationsGrpc integrationsGrpc.IntegrationsClient
	TokensGrpc       tokens.TokensClient
	BotsGrpc         bots.BotsClient
	ParserGrpc       parser.ParserClient
	SchedulerGrpc    scheduler.SchedulerClient
	TimersGrpc       timers.TimersClient
	EventSubGrpc     eventsub.EventSubClient
	DiscordGrpc      discord.DiscordClient

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
			Discord:      opts.DiscordGrpc,
		},
		Logger: opts.Logger,
	}

	return &UnProtected{
		Twitch: &twitch.Twitch{
			Deps: d,
		},
		Stats: stats.New(d),
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
				"user:read:follows",
				"channel:bot",
				"channel:manage:raids",
			},
		},
		Modules: &modules.Modules{
			Deps: d,
		},
		Songs: &songs.Songs{
			Deps: d,
		},
		Tts:            &tts.Tts{Deps: d},
		Community:      &community.Community{Deps: d},
		PublicSettings: &public_settings.PublicSettings{Deps: d},
	}
}
