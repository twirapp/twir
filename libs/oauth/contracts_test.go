package oauth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestPublicValidationErrorsSupportErrorsIs(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "credential key", err: (CredentialKey{}).Validate(), want: ErrInvalidCredential},
		{name: "credential", err: (Credential{}).Validate(), want: ErrInvalidCredential},
		{name: "refresh result", err: (RefreshResult{}).Validate(), want: ErrInvalidCredential},
		{name: "app key", err: (AppTokenKey{}).Validate(), want: ErrInvalidCredential},
		{name: "app token", err: (AppToken{}).Validate(), want: ErrInvalidCredential},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !errors.Is(test.err, test.want) {
				t.Fatalf("error = %v", test.err)
			}
		})
	}
}

func TestConstructorValidationErrorsSupportErrorsIs(t *testing.T) {
	if _, err := NewRefreshRuntime(nil, nil, nil, RuntimeOptions{Skew: -1}); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("runtime error = %v", err)
	}
	if _, err := NewAppTokenSource(AppTokenDependencies{}, AppTokenSourceOptions{WorkTimeout: -1}); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("app source error = %v", err)
	}
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	t.Cleanup(func() { _ = client.Close() })
	if _, err := NewRedisLocker(client, RedisLockerOptions{}); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("Redis locker error = %v", err)
	}
}

func TestOptionalTypedNilDependenciesUseDefaults(t *testing.T) {
	var clock *nilClock
	var observer *nilObserver
	runtime, err := NewRefreshRuntime(
		newMemoryCredentialStore(validExpiredCredential()),
		&recordingRefresher{},
		&recordingLocker{lease: newControlledLease(context.Background())},
		RuntimeOptions{Clock: clock, Observer: observer},
	)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.clock == nil || runtime.observer != nil {
		t.Fatal("typed nil optional dependencies were not normalized")
	}
}

func TestTypedNilLeaseReturnsTypedCoordinatorError(t *testing.T) {
	var lease *controlledLease
	runtime, err := NewRefreshRuntime(
		newMemoryCredentialStore(validExpiredCredential()),
		&recordingRefresher{},
		&recordingLocker{lease: lease},
		RuntimeOptions{},
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.Refresh(context.Background(), validExpiredCredential().Key()); !errors.Is(err, ErrCoordinator) || !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("runtime error = %v", err)
	}
	source := newTestAppTokenSource(t, newMemoryAppTokenStore(), &recordingAppFetcher{}, &recordingAppLocker{lease: lease})
	if _, err := source.Token(context.Background(), validAppTokenKey()); !errors.Is(err, ErrCoordinator) || !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("app source error = %v", err)
	}
}

type nilClock struct{}

func (*nilClock) Now() time.Time { panic("typed nil clock used") }

type nilObserver struct{}

func (*nilObserver) Observe(context.Context, Event) error { panic("typed nil observer used") }
