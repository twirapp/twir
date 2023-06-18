package impl

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Api struct {
	redis *redis.Client
	db    *gorm.DB
}

type Opts struct {
	Redis *redis.Client
	DB    *gorm.DB
}

func NewApi(opts Opts) *Api {
	return &Api{
		redis: opts.Redis,
		db:    opts.DB,
	}
}
