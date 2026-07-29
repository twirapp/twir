package oauth

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestAppTokenSourceCallerCancellationDoesNotPoisonOwnerOrWaiter(t *testing.T) {
	store := newMemoryAppTokenStore()
	started := make(chan struct{})
	proceed := make(chan struct{})
	fetcher := &recordingAppFetcher{hook: func(ctx context.Context, _ AppTokenKey) (AppToken, error) {
		close(started)
		select {
		case <-ctx.Done():
			return AppToken{}, ctx.Err()
		case <-proceed:
			return validAppToken(), nil
		}
	}}
	source := newTestAppTokenSource(t, store, fetcher, &recordingAppLocker{lease: newControlledLease(context.Background())})
	callerContext, cancelCaller := context.WithCancel(context.Background())
	callerDone := make(chan error, 1)
	go func() {
		_, err := source.Token(callerContext, validAppTokenKey())
		callerDone <- err
	}()
	<-started
	cancelCaller()
	if err := <-callerDone; !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled caller error = %v", err)
	}
	waiterDone := make(chan error, 1)
	go func() {
		_, err := source.Token(context.Background(), validAppTokenKey())
		waiterDone <- err
	}()
	close(proceed)
	if err := <-waiterDone; err != nil {
		t.Fatalf("waiter error = %v", err)
	}
	if fetcher.callCount() != 1 {
		t.Fatalf("fetch calls = %d", fetcher.callCount())
	}
}

func TestAppTokenSourceCachesCommittedTokenWithSkew(t *testing.T) {
	store := newMemoryAppTokenStore()
	store.put(validAppTokenKey(), AppToken{AccessToken: "cached", ObtainedAt: time.Unix(2, 0), ExpiresIn: time.Hour})
	fetcher := &recordingAppFetcher{token: validAppToken()}
	source := newTestAppTokenSource(t, store, fetcher, &recordingAppLocker{lease: newControlledLease(context.Background())})

	first, err := source.Token(context.Background(), validAppTokenKey())
	if err != nil {
		t.Fatal(err)
	}
	second, err := source.Token(context.Background(), validAppTokenKey())
	if err != nil {
		t.Fatal(err)
	}
	if first.AccessToken != "cached" || second.AccessToken != "cached" || fetcher.callCount() != 0 {
		t.Fatal("valid stored token was not cached")
	}
	loads, commits := store.counts()
	if loads != 1 || commits != 0 {
		t.Fatalf("store loads/commits = %d/%d", loads, commits)
	}
}

func TestAppTokenSourceRefreshesInsideSkew(t *testing.T) {
	store := newMemoryAppTokenStore()
	store.put(validAppTokenKey(), AppToken{AccessToken: "near-expiry", ObtainedAt: time.Unix(2, 0), ExpiresIn: time.Minute})
	fetcher := &recordingAppFetcher{token: validAppToken()}
	source := newAppTokenSourceWithOptions(t, store, fetcher, &recordingAppLocker{lease: newControlledLease(context.Background())}, AppTokenSourceOptions{
		Clock: fixedClock{now: time.Unix(32, 0)}, Skew: 31 * time.Second, WorkTimeout: time.Second,
	})

	token, err := source.Token(context.Background(), validAppTokenKey())
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != validAppToken().AccessToken || fetcher.callCount() != 1 {
		t.Fatal("token inside skew was not refreshed")
	}
}

func TestAppTokenSourceCloseCancelsAndSettlesOwner(t *testing.T) {
	store := newMemoryAppTokenStore()
	started := make(chan struct{})
	cancelled := make(chan struct{})
	settle := make(chan struct{})
	fetcher := &recordingAppFetcher{hook: func(ctx context.Context, _ AppTokenKey) (AppToken, error) {
		close(started)
		<-ctx.Done()
		close(cancelled)
		<-settle
		return AppToken{}, ctx.Err()
	}}
	source := newTestAppTokenSource(t, store, fetcher, &recordingAppLocker{lease: newControlledLease(context.Background())})
	tokenDone := make(chan error, 1)
	go func() {
		_, err := source.Token(context.Background(), validAppTokenKey())
		tokenDone <- err
	}()
	<-started
	closeDone := make(chan error, 1)
	go func() { closeDone <- source.Close() }()
	<-cancelled
	select {
	case err := <-closeDone:
		t.Errorf("Close returned before owner settled: %v", err)
	default:
	}
	close(settle)
	if err := <-tokenDone; !errors.Is(err, ErrClosed) || !errors.Is(err, ErrRefresh) {
		t.Fatalf("Token error = %v", err)
	}
	if err := <-closeDone; err != nil {
		t.Fatal(err)
	}
}

func TestAppTokenSourceCloseIsConcurrentAndPreventsNewWork(t *testing.T) {
	locker := &recordingAppLocker{lease: newControlledLease(context.Background())}
	source := newTestAppTokenSource(t, newMemoryAppTokenStore(), &recordingAppFetcher{}, locker)
	var group sync.WaitGroup
	for range 16 {
		group.Add(1)
		go func() {
			defer group.Done()
			if err := source.Close(); err != nil {
				t.Errorf("Close: %v", err)
			}
		}()
	}
	group.Wait()
	if _, err := source.Token(context.Background(), validAppTokenKey()); !errors.Is(err, ErrClosed) {
		t.Fatalf("Token error = %v", err)
	}
	if locker.calls.Load() != 0 {
		t.Fatal("work started after Close")
	}
}

func newTestAppTokenSource(t *testing.T, store AppTokenStore, fetcher AppTokenFetcher, locker AppTokenLocker) *AppTokenSource {
	t.Helper()
	return newAppTokenSourceWithOptions(t, store, fetcher, locker, AppTokenSourceOptions{
		Clock: fixedClock{now: time.Unix(3, 0)}, WorkTimeout: time.Second,
	})
}

func newAppTokenSourceWithOptions(t *testing.T, store AppTokenStore, fetcher AppTokenFetcher, locker AppTokenLocker, options AppTokenSourceOptions) *AppTokenSource {
	t.Helper()
	source, err := NewAppTokenSource(AppTokenDependencies{Store: store, Fetcher: fetcher, Locker: locker}, options)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := source.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})
	return source
}
