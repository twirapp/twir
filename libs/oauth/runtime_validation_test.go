package oauth

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRefreshRuntimeRejectsTypedNilDependencies(t *testing.T) {
	credential := validExpiredCredential()
	store := newMemoryCredentialStore(credential)
	refresher := &recordingRefresher{}
	locker := &recordingLocker{lease: newControlledLease(context.Background())}
	tests := []struct {
		name      string
		store     Store
		refresher Refresher
		locker    Locker
	}{
		{name: "store", store: (*memoryCredentialStore)(nil), refresher: refresher, locker: locker},
		{name: "refresher", store: store, refresher: (*recordingRefresher)(nil), locker: locker},
		{name: "locker", store: store, refresher: refresher, locker: (*recordingLocker)(nil)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewRefreshRuntime(test.store, test.refresher, test.locker, RuntimeOptions{})
			if !errors.Is(err, ErrInvalidOption) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestRefreshRuntimeRejectsNilContextWithoutPanic(t *testing.T) {
	runtime := newTestRuntime(t, newMemoryCredentialStore(validExpiredCredential()), &recordingRefresher{}, newControlledLease(context.Background()), fixedClock{})
	var ctx *panicContext

	_, err := runtime.Refresh(ctx, validExpiredCredential().Key())
	if !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("error = %v", err)
	}
}

func TestRefreshRuntimeRejectsInvalidLoadedCredential(t *testing.T) {
	tests := []struct {
		name       string
		credential Credential
	}{
		{name: "wrong key", credential: Credential{Provider: "other", ID: "owner", AccessToken: "value", ObtainedAt: time.Unix(1, 0), ExpiresIn: time.Hour}},
		{name: "missing access", credential: Credential{Provider: "synthetic", ID: "owner", ObtainedAt: time.Unix(1, 0), ExpiresIn: time.Hour}},
		{name: "missing expiry", credential: Credential{Provider: "synthetic", ID: "owner", AccessToken: "value", ObtainedAt: time.Unix(1, 0)}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := newMemoryCredentialStore(test.credential)
			refresher := &recordingRefresher{}
			runtime := newTestRuntime(t, store, refresher, newControlledLease(context.Background()), fixedClock{now: time.Unix(2, 0)})

			_, err := runtime.Refresh(context.Background(), validExpiredCredential().Key())
			if !errors.Is(err, ErrLoad) || !errors.Is(err, ErrInvalidCredential) {
				t.Fatalf("error = %v", err)
			}
			if refresher.callCount() != 0 {
				t.Fatal("invalid loaded credential reached refresher")
			}
		})
	}
}

func TestRefreshRuntimeRejectsInvalidRefreshResultBeforeCommit(t *testing.T) {
	tests := []RefreshResult{
		{ExpiresIn: time.Hour},
		{AccessToken: "rotated"},
	}
	for _, result := range tests {
		store := newMemoryCredentialStore(validExpiredCredential())
		runtime := newTestRuntime(t, store, &recordingRefresher{result: result}, newControlledLease(context.Background()), fixedClock{now: time.Unix(2, 0)})

		_, err := runtime.Refresh(context.Background(), validExpiredCredential().Key())
		if !errors.Is(err, ErrRefresh) || !errors.Is(err, ErrInvalidCredential) {
			t.Fatalf("error = %v", err)
		}
		_, _, commits := store.snapshot()
		if commits != 0 {
			t.Fatal("invalid result was committed")
		}
	}
}

type panicContext struct{}

func (*panicContext) Deadline() (time.Time, bool) { panic("typed nil context used") }
func (*panicContext) Done() <-chan struct{}       { panic("typed nil context used") }
func (*panicContext) Err() error                  { panic("typed nil context used") }
func (*panicContext) Value(any) any               { panic("typed nil context used") }
