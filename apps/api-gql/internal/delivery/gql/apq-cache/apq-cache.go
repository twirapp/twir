package apq_cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Redis *redis.Client
}

func New(opts Opts) *APQCache {
	return &APQCache{
		redis: opts.Redis,
	}
}

type APQCache struct {
	redis *redis.Client
}

const apqPrefix = "api-apq:"
const ttl = 7 * 24 * time.Hour

func (c *APQCache) Add(ctx context.Context, key string, value interface{}) {
	c.redis.Set(ctx, apqPrefix+key, value, ttl)
}

func (c *APQCache) Get(ctx context.Context, key string) (interface{}, bool) {
	s, err := c.redis.Get(ctx, apqPrefix+key).Result()
	if err != nil {
		return struct{}{}, false
	}
	return s, true
}
