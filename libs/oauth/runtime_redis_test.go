package oauth

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestRefreshRuntimeRealRedisTwoInstancesRefreshOnce(t *testing.T) {
	fixture := newRedisFixture(t)
	store := newMemoryCredentialStore(validExpiredCredential())
	started := make(chan struct{})
	proceed := make(chan struct{})
	var startOnce sync.Once
	refresher := &recordingRefresher{hook: func(context.Context, Credential) (RefreshResult, error) {
		startOnce.Do(func() { close(started) })
		<-proceed
		return RefreshResult{AccessToken: "rotated", ExpiresIn: time.Hour}, nil
	}}
	first := newRealRuntime(t, store, refresher, fixture.locker(t))
	second := newRealRuntime(t, store, refresher, fixture.lockerFor(t, fixture.admin))
	results := make(chan error, 2)
	go func() {
		_, err := first.Refresh(context.Background(), validExpiredCredential().Key())
		results <- err
	}()
	<-started
	go func() {
		_, err := second.Refresh(context.Background(), validExpiredCredential().Key())
		results <- err
	}()
	close(proceed)
	for range 2 {
		if err := <-results; err != nil {
			t.Fatal(err)
		}
	}
	_, loads, commits := store.snapshot()
	if refresher.callCount() != 1 || loads != 2 || commits != 1 {
		t.Fatalf("calls load/refresh/commit = %d/%d/%d", loads, refresher.callCount(), commits)
	}
}

func TestRefreshRuntimeRealRedisWaiterCancellationDoesNotPoisonOwner(t *testing.T) {
	fixture := newRedisFixture(t)
	store := newMemoryCredentialStore(validExpiredCredential())
	started := make(chan struct{})
	proceed := make(chan struct{})
	refresher := &recordingRefresher{hook: func(context.Context, Credential) (RefreshResult, error) {
		close(started)
		<-proceed
		return RefreshResult{AccessToken: "rotated", ExpiresIn: time.Hour}, nil
	}}
	owner := newRealRuntime(t, store, refresher, fixture.locker(t))
	waiter := newRealRuntime(t, store, refresher, fixture.lockerFor(t, fixture.admin))
	ownerDone := make(chan error, 1)
	go func() {
		_, err := owner.Refresh(context.Background(), validExpiredCredential().Key())
		ownerDone <- err
	}()
	<-started
	waitContext, cancelWait := context.WithTimeout(context.Background(), 75*time.Millisecond)
	defer cancelWait()
	_, waitErr := waiter.Refresh(waitContext, validExpiredCredential().Key())
	if !errors.Is(waitErr, context.DeadlineExceeded) || !errors.Is(waitErr, ErrCoordinator) {
		t.Fatalf("waiter error = %v", waitErr)
	}
	close(proceed)
	if err := <-ownerDone; err != nil {
		t.Fatalf("owner error = %v", err)
	}
	if refresher.callCount() != 1 {
		t.Fatalf("refresh calls = %d", refresher.callCount())
	}
}

func TestRefreshRuntimeRealRedisSkipsNonExpiredCredential(t *testing.T) {
	fixture := newRedisFixture(t)
	credential := validExpiredCredential()
	credential.ObtainedAt = time.Unix(2, 0)
	credential.ExpiresIn = time.Hour
	store := newMemoryCredentialStore(credential)
	refresher := &recordingRefresher{}
	runtime := newRealRuntime(t, store, refresher, fixture.locker(t))

	if _, err := runtime.Refresh(context.Background(), credential.Key()); err != nil {
		t.Fatal(err)
	}
	if refresher.callCount() != 0 {
		t.Fatal("non-expired credential was refreshed")
	}
}

func TestRefreshRuntimeRealRedisPreservesOmittedRefreshAndScopes(t *testing.T) {
	fixture := newRedisFixture(t)
	credential := validExpiredCredential()
	credential.RefreshToken = "retained"
	credential.Scopes = nil
	store := newMemoryCredentialStore(credential)
	runtime := newRealRuntime(t, store, &recordingRefresher{result: RefreshResult{
		AccessToken: "rotated", ExpiresIn: time.Hour,
	}}, fixture.locker(t))

	rotated, err := runtime.Refresh(context.Background(), credential.Key())
	if err != nil {
		t.Fatal(err)
	}
	if rotated.RefreshToken != "retained" || rotated.Scopes != nil {
		t.Fatal("omitted refresh fields were not preserved")
	}
}

func newRealRuntime(t *testing.T, store Store, refresher Refresher, locker Locker) *RefreshRuntime {
	t.Helper()
	runtime, err := NewRefreshRuntime(store, refresher, locker, RuntimeOptions{Clock: fixedClock{now: time.Unix(2, 0)}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := runtime.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})
	return runtime
}
