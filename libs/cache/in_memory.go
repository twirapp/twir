package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/maypok86/otter/v2"
)

type InMemory[V any] struct {
	cache  *otter.Cache[Key, V]
	loader Loader[V]
	ttl    time.Duration
}

type InMemoryOptions struct {
	MaxSize     int
	MinCapacity int
}

func DefaultInMemoryOptions() InMemoryOptions {
	return InMemoryOptions{
		MaxSize:     10000,
		MinCapacity: 1000,
	}
}

var _ Cache[any] = (*InMemory[any])(nil)

func NewInMemory[V any](
	options InMemoryOptions,
	ttl time.Duration,
	loader Loader[V],
) (InMemory[V], error) {
	cache, err := otter.New[Key, V](
		&otter.Options[Key, V]{
			MaximumSize:      options.MaxSize,
			InitialCapacity:  options.MinCapacity,
			ExpiryCalculator: otter.ExpiryAccessing[Key, V](ttl),
		},
	)
	if err != nil {
		return InMemory[V]{}, fmt.Errorf("new: %w", err)
	}

	return InMemory[V]{
		cache:  cache,
		loader: loader,
		ttl:    ttl,
	}, nil
}

func (im InMemory[V]) Get(ctx context.Context, key Key) (V, error) {
	var value V

	loader := otter.LoaderFunc[Key, V](
		func(ctx context.Context, key Key) (V, error) {
			loadedValue, err := im.loader(ctx, key)
			if err != nil {
				if errors.Is(err, ErrNotFound) {
					return value, otter.ErrNotFound
				}

				return value, err
			}

			return loadedValue, nil
		},
	)

	cachedValue, err := im.cache.Get(ctx, key, loader)
	if err != nil {
		if errors.Is(err, otter.ErrNotFound) {
			return value, fmt.Errorf("%w: %s", ErrNotFound, err)
		}

		return value, err
	}

	return cachedValue, nil
}

func (im InMemory[V]) Set(_ context.Context, key Key, value V) error {
	im.cache.Set(key, value)
	return nil
}

func (im InMemory[V]) SetFiltered(ctx context.Context, key Key, filter Filter[V]) error {
	value, err := im.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("get: %w", err)
	}

	filteredValue := filter(value)

	return im.Set(ctx, key, filteredValue)
}

func (im InMemory[V]) Invalidate(_ context.Context, key Key) error {
	if _, invalidated := im.cache.Invalidate(key); !invalidated {
		return ErrNotFound
	}

	return nil
}
