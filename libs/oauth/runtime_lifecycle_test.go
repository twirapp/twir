package oauth

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestRefreshRuntimeCloseCancelsAndSettlesInflightWork(t *testing.T) {
	store := newMemoryCredentialStore(validExpiredCredential())
	started := make(chan struct{})
	cancelled := make(chan struct{})
	settle := make(chan struct{})
	store.loadHook = func(ctx context.Context, _ CredentialKey) error {
		close(started)
		select {
		case <-ctx.Done():
			close(cancelled)
			<-settle
			return ctx.Err()
		case <-time.After(time.Second):
			return errors.New("load was not cancelled")
		}
	}
	runtime := newTestRuntime(t, store, &recordingRefresher{}, newControlledLease(context.Background()), fixedClock{now: time.Unix(2, 0)})
	refreshDone := make(chan error, 1)
	go func() {
		_, err := runtime.Refresh(context.Background(), validExpiredCredential().Key())
		refreshDone <- err
	}()
	<-started
	closeDone := make(chan error, 1)
	go func() { closeDone <- runtime.Close() }()

	select {
	case <-cancelled:
	case err := <-closeDone:
		t.Errorf("Close returned before cancelling and settling work: %v", err)
	}
	select {
	case err := <-closeDone:
		t.Errorf("Close returned before work settled: %v", err)
	default:
	}
	close(settle)
	if err := <-refreshDone; !errors.Is(err, ErrClosed) || !errors.Is(err, ErrLoad) {
		t.Fatalf("Refresh error = %v", err)
	}
	if err := <-closeDone; err != nil {
		t.Fatalf("Close error = %v", err)
	}
}

func TestRefreshRuntimeCloseIsConcurrentAndPreventsNewWork(t *testing.T) {
	locker := &recordingLocker{lease: newControlledLease(context.Background())}
	runtime, err := NewRefreshRuntime(newMemoryCredentialStore(validExpiredCredential()), &recordingRefresher{}, locker, RuntimeOptions{})
	if err != nil {
		t.Fatal(err)
	}
	var group sync.WaitGroup
	for range 16 {
		group.Add(1)
		go func() {
			defer group.Done()
			if closeErr := runtime.Close(); closeErr != nil {
				t.Errorf("Close: %v", closeErr)
			}
		}()
	}
	group.Wait()

	_, err = runtime.Refresh(context.Background(), validExpiredCredential().Key())
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("Refresh error = %v", err)
	}
	if locker.calls.Load() != 0 {
		t.Fatal("work started after Close")
	}
}

func TestRefreshRuntimeChecksLeaseAtEveryBoundary(t *testing.T) {
	tests := []struct {
		name        string
		configure   func(*controlledLease, *memoryCredentialStore, *recordingRefresher, *stagedClock)
		wantLoads   int
		wantRefresh int
		wantCommits int
	}{
		{name: "before load", configure: func(lease *controlledLease, _ *memoryCredentialStore, _ *recordingRefresher, _ *stagedClock) {
			lease.lose(ErrLeaseLost)
		}},
		{name: "after load", wantLoads: 1, configure: func(lease *controlledLease, store *memoryCredentialStore, _ *recordingRefresher, _ *stagedClock) {
			store.loadHook = func(context.Context, CredentialKey) error { lease.lose(ErrLeaseLost); return nil }
		}},
		{name: "before refresh", wantLoads: 1, configure: func(lease *controlledLease, _ *memoryCredentialStore, _ *recordingRefresher, clock *stagedClock) {
			clock.cancelAt, clock.cancel = 2, func() { lease.lose(ErrLeaseLost) }
		}},
		{name: "after refresh", wantLoads: 1, wantRefresh: 1, configure: func(lease *controlledLease, _ *memoryCredentialStore, refresher *recordingRefresher, _ *stagedClock) {
			refresher.hook = func(context.Context, Credential) (RefreshResult, error) {
				lease.lose(ErrLeaseLost)
				return RefreshResult{AccessToken: "rotated", ExpiresIn: time.Hour}, nil
			}
		}},
		{name: "before commit", wantLoads: 1, wantRefresh: 1, configure: func(lease *controlledLease, _ *memoryCredentialStore, _ *recordingRefresher, clock *stagedClock) {
			clock.cancelAt, clock.cancel = 3, func() { lease.lose(ErrLeaseLost) }
		}},
		{name: "after commit", wantLoads: 1, wantRefresh: 1, wantCommits: 1, configure: func(lease *controlledLease, store *memoryCredentialStore, _ *recordingRefresher, _ *stagedClock) {
			store.commitHook = func(context.Context, Credential) error { lease.lose(ErrLeaseLost); return nil }
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lease := newControlledLease(context.Background())
			store := newMemoryCredentialStore(validExpiredCredential())
			refresher := &recordingRefresher{result: RefreshResult{AccessToken: "rotated", ExpiresIn: time.Hour}}
			clock := &stagedClock{now: time.Unix(2, 0)}
			test.configure(lease, store, refresher, clock)
			runtime := newTestRuntime(t, store, refresher, lease, clock)

			credential, err := runtime.Refresh(context.Background(), validExpiredCredential().Key())
			if !errors.Is(err, ErrLeaseLost) {
				t.Fatalf("error = %v", err)
			}
			if credential.AccessToken != "" {
				t.Fatal("returned credential after lease loss")
			}
			_, loads, commits := store.snapshot()
			if loads != test.wantLoads || refresher.callCount() != test.wantRefresh || commits != test.wantCommits {
				t.Fatalf("calls load/refresh/commit = %d/%d/%d", loads, refresher.callCount(), commits)
			}
		})
	}
}
