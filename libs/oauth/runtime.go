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

	mu        sync.Mutex
	closed    bool
	nextID    uint64
	active    map[uint64]context.CancelCauseFunc
	activeWg  sync.WaitGroup
	closeOnce sync.Once
}

func NewRefreshRuntime(store Store, refresher Refresher, locker Locker, options RuntimeOptions) (*RefreshRuntime, error) {
	if isNil(store) || isNil(refresher) || isNil(locker) || options.Skew < 0 {
		return nil, fmt.Errorf("%w: dependencies and skew", ErrInvalidOption)
	}
	clock := options.Clock
	if isNil(clock) {
		clock = systemClock{}
	}
	observer := options.Observer
	if isNil(observer) {
		observer = nil
	}
	return &RefreshRuntime{
		store: store, refresher: refresher, locker: locker,
		clock: clock, skew: options.Skew, observer: observer,
		active: make(map[uint64]context.CancelCauseFunc),
	}, nil
}

// Close cancels all in-flight operations, waits for their release and observer
// callbacks to settle, and prevents new work. Concurrent calls are idempotent.
func (runtime *RefreshRuntime) Close() error {
	runtime.closeOnce.Do(func() {
		runtime.mu.Lock()
		runtime.closed = true
		cancellations := make([]context.CancelCauseFunc, 0, len(runtime.active))
		for _, cancel := range runtime.active {
			cancellations = append(cancellations, cancel)
		}
		runtime.mu.Unlock()
		for _, cancel := range cancellations {
			cancel(ErrClosed)
		}
		runtime.activeWg.Wait()
	})
	return nil
}

func (runtime *RefreshRuntime) Refresh(ctx context.Context, key CredentialKey) (credential Credential, err error) {
	if isNil(ctx) {
		return Credential{}, fmt.Errorf("%w: nil context", ErrInvalidOption)
	}
	if err := key.Validate(); err != nil {
		return Credential{}, err
	}
	operationContext, finish, err := runtime.begin(ctx)
	if err != nil {
		return Credential{}, err
	}
	defer finish()
	started := runtime.clock.Now()
	defer func() { runtime.observe(ctx, key, started, err) }()

	lease, err := runtime.locker.Acquire(operationContext, key)
	if err != nil {
		return Credential{}, joinContextCause(fmt.Errorf("%w: acquire: %w", ErrCoordinator, err), operationContext)
	}
	if isNil(lease) {
		return Credential{}, errors.Join(ErrCoordinator, ErrInvalidOption)
	}
	defer func() {
		if releaseErr := releaseLease(ctx, lease); releaseErr != nil {
			credential = Credential{}
			err = errors.Join(err, releaseErr)
		}
	}()
	leaseContext := lease.Context()
	workContext, stopWork, err := combinedLeaseContext(operationContext, leaseContext)
	if err != nil {
		return Credential{}, errors.Join(ErrCoordinator, err)
	}
	defer stopWork()

	if err := checkLease(lease, workContext); err != nil {
		return Credential{}, err
	}
	loaded, err := runtime.store.Load(workContext, key)
	if err != nil {
		return Credential{}, joinContextCause(fmt.Errorf("%w: %w", ErrLoad, err), operationContext, workContext)
	}
	loaded = cloneCredential(loaded)
	if loaded.Key() != key {
		return Credential{}, errors.Join(ErrLoad, fmt.Errorf("%w: loaded credential key mismatch", ErrInvalidCredential))
	}
	if err := loaded.Validate(); err != nil {
		return Credential{}, errors.Join(ErrLoad, err)
	}
	if err := checkLease(lease, workContext); err != nil {
		return Credential{}, err
	}
	if !loaded.Expired(runtime.clock.Now(), runtime.skew) {
		return cloneCredential(loaded), nil
	}
	if err := checkLease(lease, workContext); err != nil {
		return Credential{}, err
	}
	result, err := runtime.refresher.Refresh(workContext, cloneCredential(loaded))
	if err != nil {
		return Credential{}, joinContextCause(fmt.Errorf("%w: %w", ErrRefresh, err), operationContext, workContext)
	}
	result = cloneRefreshResult(result)
	if err := checkLease(lease, workContext); err != nil {
		return Credential{}, err
	}
	if err := result.Validate(); err != nil {
		return Credential{}, errors.Join(ErrRefresh, err)
	}

	rotated := cloneCredential(loaded)
	rotated.AccessToken = result.AccessToken
	rotated.ExpiresIn = result.ExpiresIn
	rotated.ObtainedAt = runtime.clock.Now()
	if result.RefreshToken != nil {
		rotated.RefreshToken = *result.RefreshToken
	}
	if result.Scopes != nil {
		rotated.Scopes = cloneRefreshResult(result).Scopes
	}
	if err := rotated.Validate(); err != nil {
		return Credential{}, errors.Join(ErrRefresh, err)
	}
	if err := checkLease(lease, workContext); err != nil {
		return Credential{}, err
	}
	if err := runtime.store.Commit(workContext, cloneCredential(rotated)); err != nil {
		return Credential{}, joinContextCause(fmt.Errorf("%w: %w", ErrCommit, err), operationContext, workContext)
	}
	if err := checkLease(lease, workContext); err != nil {
		return Credential{}, err
	}
	return cloneCredential(rotated), nil
}

func (runtime *RefreshRuntime) begin(parent context.Context) (context.Context, func(), error) {
	runtime.mu.Lock()
	if runtime.closed {
		runtime.mu.Unlock()
		return nil, func() {}, ErrClosed
	}
	operationContext, cancel := context.WithCancelCause(parent)
	runtime.nextID++
	id := runtime.nextID
	runtime.active[id] = cancel
	runtime.activeWg.Add(1)
	runtime.mu.Unlock()
	return operationContext, func() {
		cancel(context.Canceled)
		runtime.mu.Lock()
		delete(runtime.active, id)
		runtime.mu.Unlock()
		runtime.activeWg.Done()
	}, nil
}

func (runtime *RefreshRuntime) observe(ctx context.Context, key CredentialKey, started time.Time, operationErr error) {
	if runtime.observer == nil {
		return
	}
	event := Event{Provider: key.Provider, CredentialID: key.ID, Operation: "refresh", Duration: runtime.clock.Now().Sub(started)}
	if operationErr != nil {
		event.ErrorClass = errorClass(operationErr)
	}
	defer func() { _ = recover() }()
	_ = runtime.observer.Observe(ctx, event)
}

func errorClass(err error) string {
	switch {
	case errors.Is(err, ErrClosed):
		return "closed"
	case errors.Is(err, ErrLeaseLost):
		return "lease_lost"
	case errors.Is(err, ErrLoad):
		return "load"
	case errors.Is(err, ErrRefresh):
		return "refresh"
	case errors.Is(err, ErrCommit):
		return "commit"
	case errors.Is(err, ErrCoordinator):
		return "coordinator"
	case errors.Is(err, ErrInvalidCredential):
		return "invalid_credential"
	default:
		return "other"
	}
}
