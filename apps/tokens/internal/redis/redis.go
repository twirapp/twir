package redis

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
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

func NewRedisLock(redis *redis.Client) *redsync.Redsync {
	pool := goredis.NewPool(redis)
	return redsync.New(pool)
}
