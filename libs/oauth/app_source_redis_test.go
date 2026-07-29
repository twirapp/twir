package oauth

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestAppTokenSourceRealRedisTwoSourcesFetchOnce(t *testing.T) {
	fixture := newRedisFixture(t)
	store := newMemoryAppTokenStore()
	started := make(chan struct{})
	proceed := make(chan struct{})
	var once sync.Once
	fetcher := &recordingAppFetcher{hook: func(context.Context, AppTokenKey) (AppToken, error) {
		once.Do(func() { close(started) })
		<-proceed
		return validAppToken(), nil
	}}
	first := newTestAppTokenSource(t, store, fetcher, fixture.locker(t))
	second := newTestAppTokenSource(t, store, fetcher, fixture.lockerFor(t, fixture.admin))
	results := make(chan error, 2)
	go func() {
		_, err := first.Token(context.Background(), validAppTokenKey())
		results <- err
	}()
	<-started
	go func() {
		_, err := second.Token(context.Background(), validAppTokenKey())
		results <- err
	}()
	close(proceed)
	for range 2 {
		if err := <-results; err != nil {
			t.Fatal(err)
		}
	}
	loads, commits := store.counts()
	if fetcher.callCount() != 1 || loads != 2 || commits != 1 {
		t.Fatalf("calls load/fetch/commit = %d/%d/%d", loads, fetcher.callCount(), commits)
	}
}

func TestAppTokenSourceRealRedisUsesSeparateNamespace(t *testing.T) {
	fixture := newRedisFixture(t)
	locker := fixture.locker(t)
	credentialKey := CredentialKey{Provider: validAppTokenKey().Provider, ID: validAppTokenKey().ID}
	credentialLease, err := locker.Acquire(context.Background(), credentialKey)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = credentialLease.Release(context.Background()) }()
	appLease, err := locker.AcquireAppToken(context.Background(), validAppTokenKey())
	if err != nil {
		t.Fatalf("app lock collided with user credential lock: %v", err)
	}
	if err := appLease.Release(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestAppTokenSourceRealRedisErrorsAndClose(t *testing.T) {
	fixture := newRedisFixture(t)
	cause := errors.New("synthetic failure")
	store := newMemoryAppTokenStore()
	store.commitErr = cause
	source := newTestAppTokenSource(t, store, &recordingAppFetcher{token: validAppToken()}, fixture.locker(t))
	if _, err := source.Token(context.Background(), validAppTokenKey()); !errors.Is(err, ErrCommit) || !errors.Is(err, cause) {
		t.Fatalf("commit error = %v", err)
	}

	held, err := fixture.locker(t).AcquireAppToken(context.Background(), validAppTokenKey())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = held.Release(context.Background()) }()
	signal := &signalRedisHook{called: make(chan struct{})}
	waiterLocker, err := NewRedisLocker(fixture.admin, RedisLockerOptions{
		Prefix: "oauth-acceptance", TTL: 600 * time.Millisecond,
		RenewEvery: 150 * time.Millisecond, Timeout: 100 * time.Millisecond,
		Hooks: []redis.Hook{signal},
	})
	if err != nil {
		t.Fatal(err)
	}
	waiter := newTestAppTokenSource(t, newMemoryAppTokenStore(), &recordingAppFetcher{token: validAppToken()}, waiterLocker)
	done := make(chan error, 1)
	go func() {
		_, err := waiter.Token(context.Background(), validAppTokenKey())
		done <- err
	}()
	<-signal.called
	if err := waiter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := <-done; !errors.Is(err, ErrClosed) || !errors.Is(err, ErrCoordinator) {
		t.Fatalf("closed waiter error = %v", err)
	}
}

func TestAppTokenSourceRealRedisLossDuringFetchAndCommit(t *testing.T) {
	for _, stage := range []string{"fetch", "commit"} {
		t.Run(stage, func(t *testing.T) {
			fixture := newRedisFixture(t)
			locker := fixture.locker(t)
			store := newMemoryAppTokenStore()
			fetcher := &recordingAppFetcher{token: validAppToken()}
			if stage == "fetch" {
				fetcher.hook = func(ctx context.Context, _ AppTokenKey) (AppToken, error) {
					replaceAppOwner(t, fixture.admin, locker, validAppTokenKey())
					waitForDone(t, ctx.Done(), "lease loss during app fetch")
					return validAppToken(), nil
				}
			} else {
				store.commitHook = func(ctx context.Context, _ AppTokenKey, _ AppToken) error {
					replaceAppOwner(t, fixture.admin, locker, validAppTokenKey())
					waitForDone(t, ctx.Done(), "lease loss during app commit")
					return nil
				}
			}
			source := newTestAppTokenSource(t, store, fetcher, locker)
			token, err := source.Token(context.Background(), validAppTokenKey())
			if !errors.Is(err, ErrLeaseLost) || token.AccessToken != "" {
				t.Fatalf("token returned after lease loss: %v", err)
			}
		})
	}
}

func replaceAppOwner(t *testing.T, admin *redis.Client, locker *RedisLocker, key AppTokenKey) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := admin.Set(ctx, locker.redisKey(appTokenLockNamespace, string(key.Provider), string(key.ID)), "replacement-owner", time.Second).Err()
	if err != nil {
		t.Fatal(err)
	}
}
