package impl_unprotected

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_unprotected/auth"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_unprotected/commands"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_unprotected/stats"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_unprotected/twitch"
	cfg "github.com/satont/tsuwari/libs/config"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type UnProtected struct {
	*twitch.Twitch
	*stats.Stats
	*commands.Commands
	*auth.Auth
}

type Opts struct {
	fx.In

	Redis          *redis.Client
	DB             *gorm.DB
	Config         *cfg.Config
	SessionManager *scs.SessionManager
}

func New(opts Opts) *UnProtected {
	d := &impl_deps.Deps{
		Redis:          opts.Redis,
		Db:             opts.DB,
		Config:         opts.Config,
		SessionManager: opts.SessionManager,
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
			},
		},
	}
}
