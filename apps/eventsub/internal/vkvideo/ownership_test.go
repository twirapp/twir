package vkvideo

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestOwnershipAcquireAllowsOneOwnerWhenTwoContendersRace(t *testing.T) {
	// Given
	store := newMemoryLockStore()
	first := newTestOwnership(t, store, newManualTicker())
	second := newTestOwnership(t, store, newManualTicker())
	start := make(chan struct{})
	results := make(chan acquireResult, 2)

	for _, ownership := range []*Ownership{first, second} {
		go func() {
			<-start
			lease, err := ownership.Acquire(context.Background(), "binding-1", func() {})
			results <- acquireResult{lease: lease, err: err}
		}()
	}

	// When
	close(start)
	firstResult := <-results
	secondResult := <-results

	// Then
	var winner *Lease
	for _, result := range []acquireResult{firstResult, secondResult} {
		if result.err == nil {
			if winner != nil {
				t.Fatal("both contenders acquired the lease")
			}
			winner = result.lease
		}
	}
	if winner == nil {
		t.Fatal("neither contender acquired the lease")
	}
	if got := store.ownerCount(); got != 1 {
		t.Fatalf("owner count = %d, want 1", got)
	}
	if err := winner.Release(context.Background()); err != nil {
		t.Fatalf("release winning lease: %v", err)
	}
}

func TestLeaseReleaseCannotUnlockSuccessorAfterExpiry(t *testing.T) {
	// Given
	store := newMemoryLockStore()
	first := newTestOwnership(t, store, newManualTicker())
	second := newTestOwnership(t, store, newManualTicker())
	third := newTestOwnership(t, store, newManualTicker())
	staleLease, err := first.Acquire(context.Background(), "binding-1", func() {})
	if err != nil {
		t.Fatalf("acquire stale lease: %v", err)
	}
	store.expire("binding-1")
	successor, err := second.Acquire(context.Background(), "binding-1", func() {})
	if err != nil {
		t.Fatalf("acquire successor: %v", err)
	}

	// When
	err = staleLease.Release(context.Background())

	// Then
	if !errors.Is(err, ErrLeaseLost) {
		t.Fatalf("stale release error = %v, want ErrLeaseLost", err)
	}
	if _, err := third.Acquire(context.Background(), "binding-1", func() {}); err == nil {
		t.Fatal("third contender acquired the successor's lease")
	}
	if err := successor.Release(context.Background()); err != nil {
		t.Fatalf("release successor: %v", err)
	}
}

func TestLeaseRenewalLossFencesBeforeCloseCallback(t *testing.T) {
	// Given
	store := newMemoryLockStore()
	ticker := newManualTicker()
	ownership := newTestOwnership(t, store, ticker)
	closedAfterFence := make(chan bool, 1)
	var lease *Lease
	lease, err := ownership.Acquire(context.Background(), "binding-1", func() {
		closedAfterFence <- lease.Context().Err() != nil
	})
	if err != nil {
		t.Fatalf("acquire lease: %v", err)
	}
	store.expire("binding-1")

	// When
	ticker.tick()

	// Then
	select {
	case fenced := <-closedAfterFence:
		if !fenced {
			t.Fatal("close callback ran before lease context was fenced")
		}
	case <-time.After(time.Second):
		t.Fatal("close callback was not invoked")
	}
	lease.Wait()
	if !errors.Is(context.Cause(lease.Context()), ErrLeaseLost) {
		t.Fatalf("lease context cause = %v, want ErrLeaseLost", context.Cause(lease.Context()))
	}
}

func TestLeaseReleaseStopsRenewal(t *testing.T) {
	// Given
	store := newMemoryLockStore()
	ticker := newManualTicker()
	ownership := newTestOwnership(t, store, ticker)
	lease, err := ownership.Acquire(context.Background(), "binding-1", func() {})
	if err != nil {
		t.Fatalf("acquire lease: %v", err)
	}
	ticker.tick()
	store.waitForExtension(t)

	// When
	if err := lease.Release(context.Background()); err != nil {
		t.Fatalf("release lease: %v", err)
	}
	ticker.tick()

	// Then
	if got := store.extensionCount(); got != 1 {
		t.Fatalf("extension count after release = %d, want 1", got)
	}
}

func TestLeaseSurvivesCallerContextCancellation(t *testing.T) {
	// Given
	store := newMemoryLockStore()
	ownership := newTestOwnership(t, store, newManualTicker())
	callerCtx, cancel := context.WithCancel(context.Background())
	closed := make(chan struct{})
	lease, err := ownership.Acquire(callerCtx, "binding-1", func() {
		close(closed)
	})
	if err != nil {
		t.Fatalf("acquire lease: %v", err)
	}

	// When
	cancel()

	// Then
	select {
	case <-lease.Context().Done():
		t.Fatal("lease context was canceled with its caller context")
	default:
	}
	select {
	case <-closed:
		t.Fatal("owned resource was closed with its caller context")
	default:
	}

	if err := lease.Release(context.Background()); err != nil {
		t.Fatalf("release lease: %v", err)
	}
	select {
	case <-closed:
	case <-time.After(time.Second):
		t.Fatal("owned resource was not closed on release")
	}

	successor, err := ownership.Acquire(context.Background(), "binding-1", func() {})
	if err != nil {
		t.Fatalf("reacquire released lease: %v", err)
	}
	if err := successor.Release(context.Background()); err != nil {
		t.Fatalf("release successor lease: %v", err)
	}
}

type acquireResult struct {
	lease *Lease
	err   error
}

func newTestOwnership(t *testing.T, store *memoryLockStore, ticker *manualTicker) *Ownership {
	t.Helper()

	ownership, err := newOwnership(
		LeaseConfig{Expiry: time.Minute, RenewInterval: 20 * time.Second},
		&memoryMutexFactory{store: store},
		manualTickerFactory{ticker: ticker},
	)
	if err != nil {
		t.Fatalf("create ownership: %v", err)
	}

	return ownership
}
