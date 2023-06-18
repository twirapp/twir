package impl_unprotected

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_unprotected/stats"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_unprotected/twitch"
	"gorm.io/gorm"
)

type UnProtected struct {
	*twitch.Twitch
	*stats.Stats
}

type Opts struct {
	Redis *redis.Client
	DB    *gorm.DB
}

func New(opts Opts) *UnProtected {
	return &UnProtected{}
}
