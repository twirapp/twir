package generic_cacher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	"gorm.io/gorm"
)

type LoadFn[T any] func(ctx context.Context, key string) (T, error)

type GenericCacher[T any] struct {
	db *gorm.DB
	kv kv.KV

	keyPrefix          string
	loadFn             LoadFn[T]
	ttl                time.Duration
	invalidateSignaler InvalidateSignaler
}

type Opts[T any] struct {
	KV kv.KV

	KeyPrefix          string
	LoadFn             LoadFn[T]
	Ttl                time.Duration
	InvalidateSignaler InvalidateSignaler
}

func New[T any](opts Opts[T]) *GenericCacher[T] {
	cacher := &GenericCacher[T]{
		kv:        opts.KV,
		keyPrefix: opts.KeyPrefix,
		loadFn:    opts.LoadFn,
		ttl:       opts.Ttl,
	}

	if opts.InvalidateSignaler != nil {
		cacher.invalidateSignaler = opts.InvalidateSignaler
		receiver := opts.InvalidateSignaler.Receiver()

		go func() {
			for key := range receiver {
				if err := opts.KV.Delete(context.TODO(), opts.KeyPrefix+key); err != nil {
					fmt.Printf("failed to delete key %s from cache: %v \n", key, err)
					continue
				}
			}
		}()
	}

	return cacher
}

func (c *GenericCacher[T]) Get(ctx context.Context, key string) (T, error) {
	var value T

	cacheBytes, err := c.kv.Get(ctx, c.keyPrefix+key).Bytes()
	if err != nil && !errors.Is(err, kv.ErrKeyNil) {
		return value, fmt.Errorf("failed to get commands from cache: %w", err)
	}

	if len(cacheBytes) > 0 {
		if err := json.Unmarshal(cacheBytes, &value); err != nil {
			return value, fmt.Errorf("failed to unmarshal commands: %w", err)
		}
		return value, nil
	}
	//
	// c.mu.Lock()
	// defer c.mu.Unlock()

	value, err = c.loadFn(ctx, key)
	if err != nil {
		return value, err
	}

	cacheBytes, err = json.Marshal(value)
	if err != nil {
		return value, fmt.Errorf("failed to marshal commands: %w", err)
	}

	if err := c.kv.Set(
		ctx,
		c.keyPrefix+key,
		cacheBytes,
		kvoptions.WithExpire(c.ttl),
	); err != nil {
		return value, fmt.Errorf("failed to set commands to cache: %w", err)
	}

	return value, nil
}

func (c *GenericCacher[T]) Invalidate(ctx context.Context, key string) error {
	if c.invalidateSignaler != nil {
		if err := c.invalidateSignaler.Send(key); err != nil {
			return fmt.Errorf("failed to send invalidate signal: %w", err)
		}
	} else {
		err := c.kv.Delete(ctx, c.keyPrefix+key)
		if err != nil {
			return fmt.Errorf("failed to delete commands from cache: %w", err)
		}
	}

	return nil
}

func (c *GenericCacher[T]) SetValue(ctx context.Context, key string, newValue T) error {
	// c.mu.Lock()
	// defer c.mu.Unlock()

	cacheBytes, err := json.Marshal(newValue)
	if err != nil {
		return fmt.Errorf("failed to marshal commands: %w", err)
	}

	if err := c.kv.Set(
		ctx,
		c.keyPrefix+key,
		cacheBytes,
		kvoptions.WithExpire(c.ttl),
	); err != nil {
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

	// c.mu.Lock()
	// defer c.mu.Unlock()

	newData := filterFn(data)

	cacheBytes, err := json.Marshal(newData)
	if err != nil {
		return fmt.Errorf("failed to marshal commands: %w", err)
	}

	if err := c.kv.Set(
		ctx,
		c.keyPrefix+key,
		cacheBytes,
		kvoptions.WithExpire(c.ttl),
	); err != nil {
		return fmt.Errorf("failed to set commands to cache: %w", err)
	}

	return nil
}
