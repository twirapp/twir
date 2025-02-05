package generic_cacher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
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
	mu        *redsync.Mutex
}

type Opts[T any] struct {
	Redis *redis.Client

	KeyPrefix string
	LoadFn    LoadFn[T]
	Ttl       time.Duration
}

func New[T any](opts Opts[T]) *GenericCacher[T] {
	pool := goredis.NewPool(opts.Redis)
	rs := redsync.New(pool)
	mutex := rs.NewMutex(opts.KeyPrefix + "-mutex")

	return &GenericCacher[T]{
		redis:     opts.Redis,
		keyPrefix: opts.KeyPrefix,
		loadFn:    opts.LoadFn,
		ttl:       opts.Ttl,
		mu:        mutex,
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

	c.mu.Lock()
	defer c.mu.Unlock()

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
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.redis.Del(ctx, c.keyPrefix+key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete commands from cache: %w", err)
	}

	return nil
}

func (c *GenericCacher[T]) SetValue(ctx context.Context, key string, newValue T) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cacheBytes, err := json.Marshal(newValue)
	if err != nil {
		return fmt.Errorf("failed to marshal commands: %w", err)
	}

	if err := c.redis.Set(
		ctx,
		c.keyPrefix+key,
		cacheBytes,
		c.ttl,
	).Err(); err != nil {
		return fmt.Errorf("failed to set commands to cache: %w", err)
	}

	return nil
}

func (c *GenericCacher[T]) SetValueFiltered(
	ctx context.Context,
	key string,
	filterFn func(data T) T,
) error {
	data, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	newData := filterFn(data)

	cacheBytes, err := json.Marshal(newData)
	if err != nil {
		return fmt.Errorf("failed to marshal commands: %w", err)
	}

	if err := c.redis.Set(
		ctx,
		c.keyPrefix+key,
		cacheBytes,
		c.ttl,
	).Err(); err != nil {
		return fmt.Errorf("failed to set commands to cache: %w", err)
	}

	return nil
}
