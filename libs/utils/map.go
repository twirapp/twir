package utils

import "sync"

type SyncMap[T any] struct {
	store map[string]T
	mu    sync.Mutex
}

func NewSyncMap[T any]() *SyncMap[T] {
	return &SyncMap[T]{
		store: make(map[string]T),
		mu:    sync.Mutex{},
	}
}

func (c *SyncMap[T]) Add(key string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

func (c *SyncMap[T]) Get(key string) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.store[key]
	return v, ok
}

func (c *SyncMap[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}
