package oauth

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestRefreshRuntimeRealRedisLossBeforeLoad(t *testing.T) {
	fixture := newRedisFixture(t)
	locker := fixture.locker(t)
	store := newMemoryCredentialStore(validExpiredCredential())
	runtime := newRealRuntime(t, store, &recordingRefresher{}, &lossOnAcquireLocker{
		t: t, base: locker, admin: fixture.admin,
	})

	credential, err := runtime.Refresh(context.Background(), validExpiredCredential().Key())
	if !errors.Is(err, ErrLeaseLost) || credential.AccessToken != "" {
		t.Fatalf("credential returned after loss: error = %v", err)
	}
	_, loads, _ := store.snapshot()
	if loads != 0 {
		t.Fatal("store loaded after lease loss")
	}
}

func TestRefreshRuntimeRealRedisLossDuringRefreshPreventsCommit(t *testing.T) {
	fixture := newRedisFixture(t)
	locker := fixture.locker(t)
	store := newMemoryCredentialStore(validExpiredCredential())
	refresher := &recordingRefresher{hook: func(ctx context.Context, _ Credential) (RefreshResult, error) {
		replaceCredentialOwner(t, fixture.admin, locker, validExpiredCredential().Key())
		waitForDone(t, ctx.Done(), "lease loss during refresh")
		return RefreshResult{AccessToken: "rotated", ExpiresIn: time.Hour}, nil
	}}
	runtime := newRealRuntime(t, store, refresher, locker)

	credential, err := runtime.Refresh(context.Background(), validExpiredCredential().Key())
	if !errors.Is(err, ErrLeaseLost) || credential.AccessToken != "" {
		t.Fatalf("credential returned after loss: error = %v", err)
	}
	_, _, commits := store.snapshot()
	if commits != 0 {
		t.Fatal("committed after lease loss during refresh")
	}
}

func TestRefreshRuntimeRealRedisLossDuringCommitReturnsNoCredential(t *testing.T) {
	fixture := newRedisFixture(t)
	locker := fixture.locker(t)
	store := newMemoryCredentialStore(validExpiredCredential())
	store.commitHook = func(ctx context.Context, _ Credential) error {
		replaceCredentialOwner(t, fixture.admin, locker, validExpiredCredential().Key())
		waitForDone(t, ctx.Done(), "lease loss during commit")
		return nil
	}
	runtime := newRealRuntime(t, store, &recordingRefresher{result: RefreshResult{
		AccessToken: "rotated", ExpiresIn: time.Hour,
	}}, locker)

	credential, err := runtime.Refresh(context.Background(), validExpiredCredential().Key())
	if !errors.Is(err, ErrLeaseLost) || credential.AccessToken != "" {
		t.Fatalf("credential returned after loss: error = %v", err)
	}
}

func TestRefreshRuntimeRealRedisFailureClasses(t *testing.T) {
	fixture := newRedisFixture(t)
	cause := errors.New("synthetic failure")
	tests := []struct {
		name      string
		store     *memoryCredentialStore
		refresher *recordingRefresher
		want      error
	}{
		{name: "load", store: &memoryCredentialStore{loadErr: cause}, refresher: &recordingRefresher{}, want: ErrLoad},
		{name: "refresh", store: newMemoryCredentialStore(validExpiredCredential()), refresher: &recordingRefresher{err: cause}, want: ErrRefresh},
		{name: "commit", store: &memoryCredentialStore{credential: validExpiredCredential(), commitErr: cause}, refresher: &recordingRefresher{result: RefreshResult{AccessToken: "rotated", ExpiresIn: time.Hour}}, want: ErrCommit},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runtime := newRealRuntime(t, test.store, test.refresher, fixture.locker(t))
			_, err := runtime.Refresh(context.Background(), validExpiredCredential().Key())
			if !errors.Is(err, test.want) || !errors.Is(err, cause) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestRefreshRuntimeRealRedisCloseWhileWaiting(t *testing.T) {
	fixture := newRedisFixture(t)
	key := validExpiredCredential().Key()
	held, err := fixture.locker(t).Acquire(context.Background(), key)
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
	runtime := newRealRuntime(t, newMemoryCredentialStore(validExpiredCredential()), &recordingRefresher{}, waiterLocker)
	refreshDone := make(chan error, 1)
	go func() {
		_, err := runtime.Refresh(context.Background(), key)
		refreshDone <- err
	}()
	<-signal.called
	if err := runtime.Close(); err != nil {
		t.Fatal(err)
	}
	if err := <-refreshDone; !errors.Is(err, ErrClosed) || !errors.Is(err, ErrCoordinator) {
		t.Fatalf("Refresh error = %v", err)
	}
}

type lossOnAcquireLocker struct {
	t     *testing.T
	base  *RedisLocker
	admin *redis.Client
}

func (locker *lossOnAcquireLocker) Acquire(ctx context.Context, key CredentialKey) (Lease, error) {
	lease, err := locker.base.Acquire(ctx, key)
	if err != nil {
		return nil, err
	}
	replaceCredentialOwner(locker.t, locker.admin, locker.base, key)
	waitForDone(locker.t, lease.Context().Done(), "lease loss after acquire")
	return lease, nil
}

type signalRedisHook struct {
	called chan struct{}
	once   sync.Once
}

func (*signalRedisHook) DialHook(next redis.DialHook) redis.DialHook { return next }
func (hook *signalRedisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, command redis.Cmder) error {
		hook.once.Do(func() { close(hook.called) })
		return next(ctx, command)
	}
}
func (*signalRedisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}

func replaceCredentialOwner(t *testing.T, admin *redis.Client, locker *RedisLocker, key CredentialKey) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := admin.Set(ctx, locker.redisKey(credentialLockNamespace, string(key.Provider), string(key.ID)), "replacement-owner", time.Second).Err()
	if err != nil {
		t.Fatal(err)
	}
}

func waitForDone(t *testing.T, done <-chan struct{}, description string) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal(description)
	}
}
