package impl_unprotected

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_unprotected/commands"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_unprotected/stats"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_unprotected/twitch"
	"gorm.io/gorm"
)

type UnProtected struct {
	*twitch.Twitch
	*stats.Stats
	*commands.Commands
}

type Opts struct {
	Redis *redis.Client
	DB    *gorm.DB
}

func New(opts Opts) *UnProtected {
	d := &impl_deps.Deps{
		Redis: opts.Redis,
		Db:    opts.DB,
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
	}
}
