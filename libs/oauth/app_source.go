package oauth

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type AppTokenKey string
type AppToken struct {
	AccessToken string
	ObtainedAt  time.Time
	ExpiresIn   time.Duration
}

func (t AppToken) Expired(now time.Time, skew time.Duration) bool {
	return t.AccessToken == "" || t.ObtainedAt.IsZero() || t.ExpiresIn <= 0 || !now.Before(t.ObtainedAt.Add(t.ExpiresIn).Add(-skew))
}

type AppTokenFetcher interface {
	FetchAppToken(context.Context, AppTokenKey) (AppToken, error)
}
type AppTokenSource struct {
	fetcher AppTokenFetcher
	clock   Clock
	skew    time.Duration
	mu      sync.Mutex
	cache   map[AppTokenKey]AppToken
	calls   map[AppTokenKey]*appCall
	closed  bool
}
type appCall struct {
	done  chan struct{}
	token AppToken
	err   error
}

func NewAppTokenSource(fetcher AppTokenFetcher, clock Clock, skew time.Duration) (*AppTokenSource, error) {
	if fetcher == nil || skew < 0 {
		return nil, fmt.Errorf("%w: app source", ErrInvalidOption)
	}
	if clock == nil {
		clock = systemClock{}
	}
	return &AppTokenSource{fetcher: fetcher, clock: clock, skew: skew, cache: map[AppTokenKey]AppToken{}, calls: map[AppTokenKey]*appCall{}}, nil
}
func (s *AppTokenSource) Close() error { s.mu.Lock(); defer s.mu.Unlock(); s.closed = true; return nil }
func (s *AppTokenSource) Token(ctx context.Context, key AppTokenKey) (AppToken, error) {
	if key == "" {
		return AppToken{}, fmt.Errorf("%w: app token key", ErrInvalidCredential)
	}
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return AppToken{}, ErrClosed
	}
	if token := s.cache[key]; !token.Expired(s.clock.Now(), s.skew) {
		s.mu.Unlock()
		return token, nil
	}
	if call := s.calls[key]; call != nil {
		s.mu.Unlock()
		select {
		case <-ctx.Done():
			return AppToken{}, ctx.Err()
		case <-call.done:
			return call.token, call.err
		}
	}
	call := &appCall{done: make(chan struct{})}
	s.calls[key] = call
	s.mu.Unlock()
	token, err := s.fetcher.FetchAppToken(ctx, key)
	if err == nil && token.Expired(s.clock.Now(), 0) {
		err = fmt.Errorf("%w: fetched app token", ErrInvalidCredential)
	}
	s.mu.Lock()
	call.token, call.err = token, err
	if err == nil {
		s.cache[key] = token
	}
	delete(s.calls, key)
	close(call.done)
	s.mu.Unlock()
	return token, err
}
