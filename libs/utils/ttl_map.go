package utils

import (
	"sync"
	"time"
)

type storedValue[T any] struct {
	val       T
	expiresAt time.Time
}

type TtlSyncMap[T any] struct {
	store map[string]storedValue[T]
	rLock *sync.Mutex
	wLock *sync.Mutex
	ttl   time.Duration
}

func NewTtlSyncMap[T any](ttl time.Duration) *TtlSyncMap[T] {
	return &TtlSyncMap[T]{
		store: make(map[string]storedValue[T]),
		rLock: &sync.Mutex{},
		wLock: &sync.Mutex{},
		ttl:   ttl,
	}
}

func (c *TtlSyncMap[T]) Add(key string, value T) {
	c.wLock.Lock()
	defer c.wLock.Unlock()
	c.store[key] = storedValue[T]{
		val:       value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func (c *TtlSyncMap[T]) Get(key string) (T, bool) {
	c.rLock.Lock()
	defer c.rLock.Unlock()

	v, ok := c.store[key]

	if !ok {
		return v.val, false
	}

	if time.Now().After(v.expiresAt) {
		c.Delete(key)
		return v.val, false
	}

	return v.val, true
}

func (c *TtlSyncMap[T]) Delete(key string) {
	c.wLock.Lock()
	defer c.wLock.Unlock()
	delete(c.store, key)
}

func (c *TtlSyncMap[T]) GetAll() []T {
	var result []T
	for _, v := range c.store {
		result = append(result, v.val)
	}

	return result
}
