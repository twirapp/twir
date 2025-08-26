package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

type Redis[V any] struct {
	client *redis.Client
	loader Loader[V]
	prefix string
	ttl    time.Duration
}

var _ Cache[any] = (*Redis[any])(nil)

func NewRedis[V any](
	client *redis.Client,
	prefix string,
	ttl time.Duration,
	loader Loader[V],
) Redis[V] {
	return Redis[V]{
		client: client,
		loader: loader,
		prefix: prefix,
		ttl:    ttl,
	}
}

func (r Redis[V]) Get(ctx context.Context, key Key) (V, error) {
	var (
		value    V
		valueKey = r.prefix + key
	)

	valueBytes, err := r.client.Get(ctx, valueKey).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return value, fmt.Errorf("get value: %w", err)
	}

	if len(valueBytes) > 0 {
		if err = json.Unmarshal(valueBytes, &value); err != nil {
			return value, fmt.Errorf("unmarshal value: %w", err)
		}

		return value, nil
	}

	loadedValue, err := r.loader(ctx, key)
	if err != nil {
		return value, fmt.Errorf("load key: %w", err)
	}

	if err = r.Set(ctx, key, loadedValue); err != nil {
		return value, fmt.Errorf("set loaded value: %w", err)
	}

	return loadedValue, nil
}

func (r Redis[V]) Set(ctx context.Context, key Key, value V) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal value: %w", err)
	}

	if err = r.client.Set(ctx, r.prefix+key, valueBytes, r.ttl).Err(); err != nil {
		return fmt.Errorf("set value: %w", err)
	}

	return nil
}

func (r Redis[V]) SetFiltered(ctx context.Context, key Key, filter Filter[V]) error {
	value, err := r.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}

	filteredValue := filter(value)

	if err = r.Set(ctx, key, filteredValue); err != nil {
		return fmt.Errorf("set filtered value: %w", err)
	}

	return nil
}

func (r Redis[V]) Invalidate(ctx context.Context, key Key) error {
	if err := r.client.Del(ctx, r.prefix+key).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}

		return fmt.Errorf("delete key: %w", err)
	}

	return nil
}
