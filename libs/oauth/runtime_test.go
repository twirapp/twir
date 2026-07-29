package oauth

import (
	"context"
	"testing"
	"time"
)

func TestCredentialExpiredUsesSkew(t *testing.T) {
	credential := Credential{ObtainedAt: time.Unix(0, 0), ExpiresIn: time.Hour}
	if !credential.Expired(time.Unix(0, 0).Add(50*time.Minute), 10*time.Minute) {
		t.Fatal("expected credential to be expired within skew")
	}
}

func TestRefreshRuntimePreservesOmittedRefreshToken(t *testing.T) {
	store := &memoryStore{credential: Credential{Provider: Provider("kick"), ID: CredentialID("one"), AccessToken: "old", RefreshToken: "keep", ObtainedAt: time.Unix(0, 0), ExpiresIn: time.Second}}
	runtime, err := NewRefreshRuntime(store, refreshFunc(func(context.Context, Credential) (RefreshResult, error) {
		return RefreshResult{AccessToken: "new", ExpiresIn: time.Hour}, nil
	}), immediateLocker{}, RuntimeOptions{Clock: fixedClock{now: time.Unix(2, 0)}})
	if err != nil {
		t.Fatal(err)
	}
	credential, err := runtime.Refresh(context.Background(), store.credential.Key())
	if err != nil {
		t.Fatal(err)
	}
	if credential.RefreshToken != "keep" {
		t.Fatalf("refresh token = %q", credential.RefreshToken)
	}
}

type memoryStore struct{ credential Credential }

func (s *memoryStore) Load(context.Context, CredentialKey) (Credential, error) {
	return s.credential, nil
}
func (s *memoryStore) Commit(_ context.Context, c Credential) error { s.credential = c; return nil }

type refreshFunc func(context.Context, Credential) (RefreshResult, error)

func (f refreshFunc) Refresh(ctx context.Context, c Credential) (RefreshResult, error) {
	return f(ctx, c)
}

type immediateLocker struct{}

func (immediateLocker) Acquire(ctx context.Context, _ CredentialKey) (Lease, error) {
	return immediateLease{ctx}, nil
}

type immediateLease struct{ ctx context.Context }

func (l immediateLease) Context() context.Context    { return l.ctx }
func (immediateLease) Lost() <-chan struct{}         { return nil }
func (immediateLease) Release(context.Context) error { return nil }

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }
