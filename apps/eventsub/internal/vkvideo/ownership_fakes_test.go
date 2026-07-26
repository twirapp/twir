package vkvideo

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-redsync/redsync/v4"
)

type memoryLockStore struct {
	mu             sync.Mutex
	nextToken      uint64
	owners         map[string]uint64
	extensions     int
	extensionCalls chan struct{}
}

func newMemoryLockStore() *memoryLockStore {
	return &memoryLockStore{
		owners:         make(map[string]uint64),
		extensionCalls: make(chan struct{}, 1),
	}
}

func (s *memoryLockStore) newMutex(key string) *memoryMutex {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextToken++
	return &memoryMutex{store: s, key: key, token: s.nextToken}
}

func (s *memoryLockStore) expire(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.owners, key)
}

func (s *memoryLockStore) ownerCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.owners)
}

func (s *memoryLockStore) extensionCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.extensions
}

func (s *memoryLockStore) waitForExtension(t *testing.T) {
	t.Helper()

	select {
	case <-s.extensionCalls:
	case <-time.After(time.Second):
		t.Fatal("renewal was not attempted")
	}
}

type memoryMutexFactory struct {
	store *memoryLockStore
}

func (f *memoryMutexFactory) NewMutex(key string, _ time.Duration) leaseMutex {
	return f.store.newMutex(key)
}

type memoryMutex struct {
	store *memoryLockStore
	key   string
	token uint64
}

func (m *memoryMutex) TryLockContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	m.store.mu.Lock()
	defer m.store.mu.Unlock()
	if _, exists := m.store.owners[m.key]; exists {
		return redsync.ErrFailed
	}
	m.store.owners[m.key] = m.token
	return nil
}

func (m *memoryMutex) ExtendContext(ctx context.Context) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	m.store.mu.Lock()
	m.store.extensions++
	extended := m.store.owners[m.key] == m.token
	m.store.mu.Unlock()

	select {
	case m.store.extensionCalls <- struct{}{}:
	default:
	}

	return extended, nil
}

func (m *memoryMutex) UnlockContext(ctx context.Context) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	m.store.mu.Lock()
	defer m.store.mu.Unlock()
	if m.store.owners[m.key] != m.token {
		return false, nil
	}
	delete(m.store.owners, m.key)
	return true, nil
}

type manualTickerFactory struct {
	ticker *manualTicker
}

func (f manualTickerFactory) NewTicker(time.Duration) renewalTicker {
	return f.ticker
}

type manualTicker struct {
	ticks   chan time.Time
	stopped chan struct{}
	stop    sync.Once
}

func newManualTicker() *manualTicker {
	return &manualTicker{
		ticks:   make(chan time.Time, 1),
		stopped: make(chan struct{}),
	}
}

func (t *manualTicker) Chan() <-chan time.Time {
	return t.ticks
}

func (t *manualTicker) Stop() {
	t.stop.Do(func() {
		close(t.stopped)
	})
}

func (t *manualTicker) tick() {
	t.ticks <- time.Time{}
}
