package oauth

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type memoryAppTokenStore struct {
	mu         sync.Mutex
	tokens     map[AppTokenKey]AppToken
	loadErr    error
	commitErr  error
	loadHook   func(context.Context, AppTokenKey) error
	commitHook func(context.Context, AppTokenKey, AppToken) error
	loads      int
	commits    int
}

func newMemoryAppTokenStore() *memoryAppTokenStore {
	return &memoryAppTokenStore{tokens: make(map[AppTokenKey]AppToken)}
}

func (store *memoryAppTokenStore) LoadAppToken(ctx context.Context, key AppTokenKey) (AppToken, error) {
	store.mu.Lock()
	store.loads++
	token, found := store.tokens[key]
	err := store.loadErr
	hook := store.loadHook
	store.mu.Unlock()
	if hook != nil {
		if hookErr := hook(ctx, key); hookErr != nil {
			return AppToken{}, hookErr
		}
	}
	if err != nil {
		return AppToken{}, err
	}
	if !found {
		return AppToken{}, ErrAppTokenNotFound
	}
	return token, nil
}

func (store *memoryAppTokenStore) CommitAppToken(ctx context.Context, key AppTokenKey, token AppToken) error {
	store.mu.Lock()
	store.commits++
	err := store.commitErr
	hook := store.commitHook
	store.mu.Unlock()
	if hook != nil {
		if hookErr := hook(ctx, key, token); hookErr != nil {
			return hookErr
		}
	}
	if err != nil {
		return err
	}
	store.mu.Lock()
	store.tokens[key] = token
	store.mu.Unlock()
	return nil
}

func (store *memoryAppTokenStore) put(key AppTokenKey, token AppToken) {
	store.mu.Lock()
	store.tokens[key] = token
	store.mu.Unlock()
}

func (store *memoryAppTokenStore) counts() (int, int) {
	store.mu.Lock()
	defer store.mu.Unlock()
	return store.loads, store.commits
}

type recordingAppFetcher struct {
	mu    sync.Mutex
	token AppToken
	err   error
	hook  func(context.Context, AppTokenKey) (AppToken, error)
	calls int
}

func (fetcher *recordingAppFetcher) FetchAppToken(ctx context.Context, key AppTokenKey) (AppToken, error) {
	fetcher.mu.Lock()
	fetcher.calls++
	token := fetcher.token
	err := fetcher.err
	hook := fetcher.hook
	fetcher.mu.Unlock()
	if hook != nil {
		return hook(ctx, key)
	}
	return token, err
}

func (fetcher *recordingAppFetcher) callCount() int {
	fetcher.mu.Lock()
	defer fetcher.mu.Unlock()
	return fetcher.calls
}

type recordingAppLocker struct {
	lease Lease
	err   error
	calls atomic.Int64
}

func (locker *recordingAppLocker) AcquireAppToken(context.Context, AppTokenKey) (Lease, error) {
	locker.calls.Add(1)
	return locker.lease, locker.err
}

func validAppTokenKey() AppTokenKey {
	return AppTokenKey{Provider: Provider("synthetic"), ID: CredentialID("client")}
}

func validAppToken() AppToken {
	return AppToken{AccessToken: "app-access", ObtainedAt: time.Unix(2, 0), ExpiresIn: time.Hour}
}
