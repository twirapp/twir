package oauth

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRefreshRuntimeErrorsPreserveTypedCauses(t *testing.T) {
	loadCause := errors.New("load cause")
	refreshCause := errors.New("refresh cause")
	commitCause := errors.New("commit cause")
	coordinatorCause := errors.New("coordinator cause")
	tests := []struct {
		name      string
		store     *memoryCredentialStore
		refresher *recordingRefresher
		locker    *recordingLocker
		want      []error
	}{
		{name: "load", store: &memoryCredentialStore{loadErr: loadCause}, refresher: &recordingRefresher{}, locker: &recordingLocker{lease: newControlledLease(context.Background())}, want: []error{ErrLoad, loadCause}},
		{name: "refresh", store: newMemoryCredentialStore(validExpiredCredential()), refresher: &recordingRefresher{err: refreshCause}, locker: &recordingLocker{lease: newControlledLease(context.Background())}, want: []error{ErrRefresh, refreshCause}},
		{name: "commit", store: &memoryCredentialStore{credential: validExpiredCredential(), commitErr: commitCause}, refresher: &recordingRefresher{result: RefreshResult{AccessToken: "rotated", ExpiresIn: time.Hour}}, locker: &recordingLocker{lease: newControlledLease(context.Background())}, want: []error{ErrCommit, commitCause}},
		{name: "coordinator", store: newMemoryCredentialStore(validExpiredCredential()), refresher: &recordingRefresher{}, locker: &recordingLocker{err: coordinatorCause}, want: []error{ErrCoordinator, coordinatorCause}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runtime, err := NewRefreshRuntime(test.store, test.refresher, test.locker, RuntimeOptions{Clock: fixedClock{now: time.Unix(2, 0)}})
			if err != nil {
				t.Fatal(err)
			}
			_, err = runtime.Refresh(context.Background(), validExpiredCredential().Key())
			for _, want := range test.want {
				if !errors.Is(err, want) {
					t.Fatalf("error = %v, want cause %v", err, want)
				}
			}
		})
	}
}

func TestRefreshRuntimeJoinsOperationAndReleaseErrors(t *testing.T) {
	loadCause := errors.New("load cause")
	lease := newControlledLease(context.Background())
	lease.releaseErr = ErrLeaseLost
	store := newMemoryCredentialStore(validExpiredCredential())
	store.loadErr = loadCause
	runtime := newTestRuntime(t, store, &recordingRefresher{}, lease, fixedClock{now: time.Unix(2, 0)})

	credential, err := runtime.Refresh(context.Background(), validExpiredCredential().Key())
	for _, want := range []error{ErrLoad, loadCause, ErrCoordinator, ErrLeaseLost} {
		if !errors.Is(err, want) {
			t.Fatalf("error = %v, want cause %v", err, want)
		}
	}
	if credential.AccessToken != "" {
		t.Fatal("returned credential after release failure")
	}
}

func TestRefreshRuntimeObserverFailuresAreIsolated(t *testing.T) {
	for _, observer := range []Observer{
		observerFunc(func(context.Context, Event) error { return errors.New("observer failure") }),
		observerFunc(func(context.Context, Event) error { panic("observer panic") }),
	} {
		store := newMemoryCredentialStore(validExpiredCredential())
		lease := newControlledLease(context.Background())
		runtime, err := NewRefreshRuntime(store, &recordingRefresher{result: RefreshResult{AccessToken: "rotated", ExpiresIn: time.Hour}}, &recordingLocker{lease: lease}, RuntimeOptions{Clock: fixedClock{now: time.Unix(2, 0)}, Observer: observer})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := runtime.Refresh(context.Background(), validExpiredCredential().Key()); err != nil {
			t.Fatalf("observer altered refresh result: %v", err)
		}
	}
}

func TestObserverEventCannotContainCredentialMaterial(t *testing.T) {
	typeOfEvent := reflect.TypeOf(Event{})
	for index := 0; index < typeOfEvent.NumField(); index++ {
		name := strings.ToLower(typeOfEvent.Field(index).Name)
		if strings.Contains(name, "token") || strings.Contains(name, "scope") {
			t.Fatalf("observer event exposes credential field %q", name)
		}
	}
}

type observerFunc func(context.Context, Event) error

func (f observerFunc) Observe(ctx context.Context, event Event) error { return f(ctx, event) }

type eventRecorder struct {
	mu     sync.Mutex
	events []Event
}

func (r *eventRecorder) Observe(_ context.Context, event Event) error {
	r.mu.Lock()
	r.events = append(r.events, event)
	r.mu.Unlock()
	return nil
}
