package oauth

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCredentialExpiredUsesSkew(t *testing.T) {
	credential := Credential{ObtainedAt: time.Unix(0, 0), ExpiresIn: time.Hour}
	if !credential.Expired(time.Unix(0, 0).Add(50*time.Minute), 10*time.Minute) {
		t.Fatal("expected credential to be expired within skew")
	}
}

func TestRefreshRuntimePreservesOmittedRefreshFields(t *testing.T) {
	credential := validExpiredCredential()
	credential.RefreshToken = "retained"
	credential.Scopes = nil
	store := newMemoryCredentialStore(credential)
	runtime := newTestRuntime(t, store, &recordingRefresher{result: RefreshResult{
		AccessToken: "rotated",
		ExpiresIn:   time.Hour,
	}}, newControlledLease(context.Background()), fixedClock{now: time.Unix(2, 0)})

	rotated, err := runtime.Refresh(context.Background(), credential.Key())
	if err != nil {
		t.Fatal(err)
	}
	if rotated.RefreshToken != "retained" {
		t.Fatalf("refresh token was not preserved")
	}
	if rotated.Scopes != nil {
		t.Fatalf("nil scopes became non-nil: %#v", rotated.Scopes)
	}
}

func TestRefreshRuntimeDoesNotAliasCredentialSlices(t *testing.T) {
	credential := validExpiredCredential()
	credential.Scopes = []string{"input"}
	store := newMemoryCredentialStore(credential)
	resultScopes := []string{"result"}
	refresher := &recordingRefresher{hook: func(_ context.Context, input Credential) (RefreshResult, error) {
		input.Scopes[0] = "mutated-input"
		return RefreshResult{AccessToken: "rotated", Scopes: resultScopes, ExpiresIn: time.Hour}, nil
	}}
	runtime := newTestRuntime(t, store, refresher, newControlledLease(context.Background()), fixedClock{now: time.Unix(2, 0)})

	rotated, err := runtime.Refresh(context.Background(), credential.Key())
	if err != nil {
		t.Fatal(err)
	}
	resultScopes[0] = "mutated-result"
	rotated.Scopes[0] = "mutated-return"
	committed, _, _ := store.snapshot()
	if committed.Scopes[0] != "result" {
		t.Fatalf("committed scopes aliased external slice: %#v", committed.Scopes)
	}
	if credential.Scopes[0] != "input" {
		t.Fatalf("input scopes were mutated: %#v", credential.Scopes)
	}
}

func TestRefreshRuntimeSkipsRefreshForNonExpiredCredential(t *testing.T) {
	credential := validExpiredCredential()
	credential.ObtainedAt = time.Unix(2, 0)
	credential.ExpiresIn = time.Hour
	store := newMemoryCredentialStore(credential)
	refresher := &recordingRefresher{}
	runtime := newTestRuntime(t, store, refresher, newControlledLease(context.Background()), fixedClock{now: time.Unix(3, 0)})

	loaded, err := runtime.Refresh(context.Background(), credential.Key())
	if err != nil {
		t.Fatal(err)
	}
	if loaded.AccessToken != credential.AccessToken || refresher.callCount() != 0 {
		t.Fatal("non-expired credential was refreshed")
	}
}

func validExpiredCredential() Credential {
	return Credential{
		Provider:     Provider("synthetic"),
		ID:           CredentialID("owner"),
		AccessToken:  "access-value",
		RefreshToken: "refresh-value",
		ObtainedAt:   time.Unix(0, 0),
		ExpiresIn:    time.Second,
	}
}

func newTestRuntime(t *testing.T, store Store, refresher Refresher, lease Lease, clock Clock) *RefreshRuntime {
	t.Helper()
	runtime, err := NewRefreshRuntime(store, refresher, &recordingLocker{lease: lease}, RuntimeOptions{Clock: clock})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := runtime.Close(); err != nil && !errors.Is(err, ErrClosed) {
			t.Errorf("Close: %v", err)
		}
	})
	return runtime
}
