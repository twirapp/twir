package oauth

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type RuntimeOptions struct {
	Clock    Clock
	Skew     time.Duration
	Observer Observer
}
type RefreshRuntime struct {
	store     Store
	refresher Refresher
	locker    Locker
	clock     Clock
	skew      time.Duration
	observer  Observer
	closed    chan struct{}
	closeOnce sync.Once
}

func NewRefreshRuntime(store Store, refresher Refresher, locker Locker, options RuntimeOptions) (*RefreshRuntime, error) {
	if store == nil || refresher == nil || locker == nil || options.Skew < 0 {
		return nil, fmt.Errorf("%w: dependencies and skew", ErrInvalidOption)
	}
	clock := options.Clock
	if clock == nil {
		clock = systemClock{}
	}
	return &RefreshRuntime{store: store, refresher: refresher, locker: locker, clock: clock, skew: options.Skew, observer: options.Observer, closed: make(chan struct{})}, nil
}

func (r *RefreshRuntime) Close() error {
	r.closeOnce.Do(func() { close(r.closed) })
	return nil
}

func (r *RefreshRuntime) Refresh(ctx context.Context, key CredentialKey) (credential Credential, err error) {
	if err := key.Validate(); err != nil {
		return Credential{}, err
	}
	select {
	case <-r.closed:
		return Credential{}, ErrClosed
	default:
	}
	started := r.clock.Now()
	defer func() {
		if r.observer != nil {
			event := Event{Provider: key.Provider, CredentialID: key.ID, Operation: "refresh", Duration: r.clock.Now().Sub(started)}
			if err != nil {
				event.ErrorClass = fmt.Sprintf("%T", err)
			}
			r.observer.Observe(ctx, event)
		}
	}()
	lease, err := r.locker.Acquire(ctx, key)
	if err != nil {
		return Credential{}, fmt.Errorf("%w: acquire: %w", ErrCoordinator, err)
	}
	defer func() {
		releaseCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err = errors.Join(err, lease.Release(releaseCtx))
	}()
	credential, err = r.store.Load(lease.Context(), key)
	if err != nil {
		return Credential{}, fmt.Errorf("%w: %w", ErrLoad, err)
	}
	if !credential.Expired(r.clock.Now(), r.skew) {
		return credential, nil
	}
	result, err := r.refresher.Refresh(lease.Context(), credential)
	if err != nil {
		return Credential{}, fmt.Errorf("%w: %w", ErrRefresh, err)
	}
	select {
	case <-lease.Context().Done():
		return Credential{}, errors.Join(ErrLeaseLost, lease.Context().Err())
	default:
	}
	credential.AccessToken, credential.ExpiresIn, credential.ObtainedAt = result.AccessToken, result.ExpiresIn, r.clock.Now()
	if result.RefreshToken != nil {
		credential.RefreshToken = *result.RefreshToken
	}
	if result.Scopes != nil {
		credential.Scopes = append([]string(nil), result.Scopes...)
	}
	if err = r.store.Commit(lease.Context(), credential); err != nil {
		return Credential{}, fmt.Errorf("%w: %w", ErrCommit, err)
	}
	select {
	case <-lease.Context().Done():
		return Credential{}, errors.Join(ErrLeaseLost, lease.Context().Err())
	default:
	}
	return credential, nil
}
