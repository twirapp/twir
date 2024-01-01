package redis

import (
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
)

func New(cfg config.Config) (*redis.Client, error) {
	params, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(params)

	return rdb, nil
}
