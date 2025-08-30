package cache

import (
	"context"
	"errors"
)

// ErrNotFound indicates that entry with key on which some operation was supposed (get, invalidate etc.)
// does not exist in cache or loader source if it's implemented.
var ErrNotFound = errors.New("entry does not exists")

type Key = string

type (
	Filter[V any] func(V) V
	Loader[V any] func(ctx context.Context, key Key) (V, error)
)

type Cache[V any] interface {
	Get(ctx context.Context, key Key) (V, error)
	Set(ctx context.Context, key Key, value V) error
	SetFiltered(ctx context.Context, key Key, filter Filter[V]) error
	Invalidate(ctx context.Context, key Key) error
}
