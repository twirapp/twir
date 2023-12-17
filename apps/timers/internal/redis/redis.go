package redis

import (
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
)

func New(config cfg.Config) (*redis.Client, error) {
	redisOpts, err := redis.ParseURL(config.RedisUrl)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(redisOpts), nil
}
