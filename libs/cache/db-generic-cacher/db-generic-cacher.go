package db_generic_cacher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type LoadFn[T any] func(ctx context.Context, key string) (T, error)

type GenericCacher[T any] struct {
	db    *gorm.DB
	redis *redis.Client

	keyPrefix string
	loadFn    LoadFn[T]
	ttl       time.Duration
}

type Opts[T any] struct {
	Redis *redis.Client

	KeyPrefix string
	LoadFn    LoadFn[T]
	Ttl       time.Duration
}

func New[T any](opts Opts[T]) *GenericCacher[T] {
	return &GenericCacher[T]{
		redis:     opts.Redis,
		keyPrefix: opts.KeyPrefix,
		loadFn:    opts.LoadFn,
		ttl:       opts.Ttl,
	}
}

func (c *GenericCacher[T]) Get(ctx context.Context, key string) (T, error) {
	var value T

	cacheBytes, err := c.redis.Get(ctx, c.keyPrefix+key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return value, fmt.Errorf("failed to get commands from cache: %w", err)
	}

	if len(cacheBytes) > 0 {
		if err := json.Unmarshal(cacheBytes, &value); err != nil {
			return value, fmt.Errorf("failed to unmarshal commands: %w", err)
		}
		return value, nil
	}

	value, err = c.loadFn(ctx, key)
	if err != nil {
		return value, err
	}

	cacheBytes, err = json.Marshal(value)
	if err != nil {
		return value, fmt.Errorf("failed to marshal commands: %w", err)
	}

	if err := c.redis.Set(
		ctx,
		c.keyPrefix+key,
		cacheBytes,
		c.ttl,
	).Err(); err != nil {
		return value, fmt.Errorf("failed to set commands to cache: %w", err)
	}

	return value, nil
}

func (c *GenericCacher[T]) Invalidate(ctx context.Context, key string) error {
	err := c.redis.Del(ctx, c.keyPrefix+key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete commands from cache: %w", err)
	}

	return nil
}
