package redis

import (
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	Cfg cfg.Config
}

func New(opts Opts) (*redis.Client, error) {
	redisOpts, err := redis.ParseURL(opts.Cfg.RedisUrl)
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(redisOpts)
	return redisClient, nil
}
