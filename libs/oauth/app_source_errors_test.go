package oauth

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestAppTokenSourceRejectsTypedNilDependencies(t *testing.T) {
	store := newMemoryAppTokenStore()
	fetcher := &recordingAppFetcher{}
	locker := &recordingAppLocker{lease: newControlledLease(context.Background())}
	tests := []AppTokenDependencies{
		{Store: (*memoryAppTokenStore)(nil), Fetcher: fetcher, Locker: locker},
		{Store: store, Fetcher: (*recordingAppFetcher)(nil), Locker: locker},
		{Store: store, Fetcher: fetcher, Locker: (*recordingAppLocker)(nil)},
	}
	for _, dependencies := range tests {
		if _, err := NewAppTokenSource(dependencies, AppTokenSourceOptions{}); !errors.Is(err, ErrInvalidOption) {
			t.Fatalf("error = %v", err)
		}
	}
}

func TestAppTokenSourceRejectsNilContextAndInvalidKey(t *testing.T) {
	source := newTestAppTokenSource(t, newMemoryAppTokenStore(), &recordingAppFetcher{}, &recordingAppLocker{lease: newControlledLease(context.Background())})
	var nilContext *panicContext
	if _, err := source.Token(nilContext, validAppTokenKey()); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("nil context error = %v", err)
	}
	if _, err := source.Token(context.Background(), AppTokenKey{}); !errors.Is(err, ErrInvalidCredential) {
		t.Fatalf("invalid key error = %v", err)
	}
}

func TestAppTokenSourcePreservesTypedOperationErrors(t *testing.T) {
	cause := errors.New("synthetic cause")
	tests := []struct {
		name    string
		store   *memoryAppTokenStore
		fetcher *recordingAppFetcher
		locker  *recordingAppLocker
		want    error
	}{
		{name: "load", store: &memoryAppTokenStore{loadErr: cause}, fetcher: &recordingAppFetcher{}, locker: &recordingAppLocker{lease: newControlledLease(context.Background())}, want: ErrLoad},
		{name: "fetch", store: newMemoryAppTokenStore(), fetcher: &recordingAppFetcher{err: cause}, locker: &recordingAppLocker{lease: newControlledLease(context.Background())}, want: ErrRefresh},
		{name: "commit", store: &memoryAppTokenStore{tokens: make(map[AppTokenKey]AppToken), commitErr: cause}, fetcher: &recordingAppFetcher{token: validAppToken()}, locker: &recordingAppLocker{lease: newControlledLease(context.Background())}, want: ErrCommit},
		{name: "coordinator", store: newMemoryAppTokenStore(), fetcher: &recordingAppFetcher{}, locker: &recordingAppLocker{err: cause}, want: ErrCoordinator},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := newTestAppTokenSource(t, test.store, test.fetcher, test.locker)
			_, err := source.Token(context.Background(), validAppTokenKey())
			if !errors.Is(err, test.want) || !errors.Is(err, cause) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestAppTokenSourceRejectsInvalidLoadedAndFetchedTokens(t *testing.T) {
	invalidTokens := []AppToken{
		{ObtainedAt: time.Unix(1, 0), ExpiresIn: time.Hour},
		{AccessToken: "value", ExpiresIn: time.Hour},
		{AccessToken: "value", ObtainedAt: time.Unix(1, 0)},
	}
	for _, token := range invalidTokens {
		store := newMemoryAppTokenStore()
		store.put(validAppTokenKey(), token)
		source := newTestAppTokenSource(t, store, &recordingAppFetcher{}, &recordingAppLocker{lease: newControlledLease(context.Background())})
		_, err := source.Token(context.Background(), validAppTokenKey())
		if !errors.Is(err, ErrLoad) || !errors.Is(err, ErrInvalidCredential) {
			t.Fatalf("loaded token error = %v", err)
		}

		fetchSource := newTestAppTokenSource(t, newMemoryAppTokenStore(), &recordingAppFetcher{token: token}, &recordingAppLocker{lease: newControlledLease(context.Background())})
		_, err = fetchSource.Token(context.Background(), validAppTokenKey())
		if !errors.Is(err, ErrRefresh) || !errors.Is(err, ErrInvalidCredential) {
			t.Fatalf("fetched token error = %v", err)
		}
	}
}
