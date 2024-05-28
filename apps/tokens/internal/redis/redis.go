package redis

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

func NewRedisLock(redis *redis.Client) *redsync.Redsync {
	pool := goredis.NewPool(redis)
	return redsync.New(pool)
}
