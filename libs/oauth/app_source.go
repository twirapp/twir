package oauth

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

const defaultAppTokenWorkTimeout = 30 * time.Second

type AppTokenSource struct {
	store       AppTokenStore
	fetcher     AppTokenFetcher
	locker      AppTokenLocker
	clock       Clock
	skew        time.Duration
	workTimeout time.Duration

	mu        sync.Mutex
	cache     map[AppTokenKey]AppToken
	calls     map[AppTokenKey]*appCall
	closed    bool
	active    sync.WaitGroup
	closeOnce sync.Once
}

type appCall struct {
	done    chan struct{}
	ctx     context.Context
	cancel  context.CancelCauseFunc
	cleanup context.CancelFunc
	token   AppToken
	err     error
}

func NewAppTokenSource(dependencies AppTokenDependencies, options AppTokenSourceOptions) (*AppTokenSource, error) {
	if isNil(dependencies.Store) || isNil(dependencies.Fetcher) || isNil(dependencies.Locker) ||
		options.Skew < 0 || options.WorkTimeout < 0 {
		return nil, fmt.Errorf("%w: app token source", ErrInvalidOption)
	}
	clock := options.Clock
	if isNil(clock) {
		clock = systemClock{}
	}
	workTimeout := options.WorkTimeout
	if workTimeout == 0 {
		workTimeout = defaultAppTokenWorkTimeout
	}
	return &AppTokenSource{
		store: dependencies.Store, fetcher: dependencies.Fetcher, locker: dependencies.Locker,
		clock: clock, skew: options.Skew, workTimeout: workTimeout,
		cache: make(map[AppTokenKey]AppToken), calls: make(map[AppTokenKey]*appCall),
	}, nil
}

// Close cancels distributed owner work, waits for it to settle, and prevents new
// token loads. Concurrent calls are idempotent.
func (source *AppTokenSource) Close() error {
	source.closeOnce.Do(func() {
		source.mu.Lock()
		source.closed = true
		cancellations := make([]context.CancelCauseFunc, 0, len(source.calls))
		for _, call := range source.calls {
			cancellations = append(cancellations, call.cancel)
		}
		source.mu.Unlock()
		for _, cancel := range cancellations {
			cancel(ErrClosed)
		}
		source.active.Wait()
	})
	return nil
}

// Token detaches newly elected owner work from the initiating caller so one
// cancelled waiter cannot poison other waiters. Work remains bounded by
// WorkTimeout and the source Close lifecycle.
func (source *AppTokenSource) Token(ctx context.Context, key AppTokenKey) (AppToken, error) {
	if isNil(ctx) {
		return AppToken{}, fmt.Errorf("%w: nil context", ErrInvalidOption)
	}
	if cause := context.Cause(ctx); cause != nil {
		return AppToken{}, cause
	}
	if err := key.Validate(); err != nil {
		return AppToken{}, err
	}
	source.mu.Lock()
	if source.closed {
		source.mu.Unlock()
		return AppToken{}, ErrClosed
	}
	if token, found := source.cache[key]; found && !token.Expired(source.clock.Now(), source.skew) {
		source.mu.Unlock()
		return token, nil
	}
	call := source.calls[key]
	if call == nil {
		timeoutContext, cleanup := context.WithTimeout(context.WithoutCancel(ctx), source.workTimeout)
		ownerContext, cancel := context.WithCancelCause(timeoutContext)
		call = &appCall{done: make(chan struct{}), ctx: ownerContext, cancel: cancel, cleanup: cleanup}
		source.calls[key] = call
		source.active.Add(1)
		go source.runCall(key, call)
	}
	source.mu.Unlock()
	select {
	case <-ctx.Done():
		return AppToken{}, context.Cause(ctx)
	case <-call.done:
		return call.token, call.err
	}
}

func (source *AppTokenSource) runCall(key AppTokenKey, call *appCall) {
	token, err := source.loadOrFetch(call.ctx, key)
	call.cleanup()
	source.mu.Lock()
	if source.closed {
		token = AppToken{}
		err = errors.Join(err, ErrClosed)
	}
	call.token = token
	call.err = err
	if err == nil {
		source.cache[key] = token
	}
	delete(source.calls, key)
	close(call.done)
	source.mu.Unlock()
	source.active.Done()
}

func (source *AppTokenSource) loadOrFetch(ctx context.Context, key AppTokenKey) (token AppToken, err error) {
	lease, err := source.locker.AcquireAppToken(ctx, key)
	if err != nil {
		return AppToken{}, joinContextCause(fmt.Errorf("%w: acquire app token: %w", ErrCoordinator, err), ctx)
	}
	if isNil(lease) {
		return AppToken{}, errors.Join(ErrCoordinator, ErrInvalidOption)
	}
	defer func() {
		if releaseErr := releaseLease(ctx, lease); releaseErr != nil {
			token = AppToken{}
			err = errors.Join(err, releaseErr)
		}
	}()
	workContext, stopWork, err := combinedLeaseContext(ctx, lease.Context())
	if err != nil {
		return AppToken{}, errors.Join(ErrCoordinator, err)
	}
	defer stopWork()
	if err := checkLease(lease, workContext); err != nil {
		return AppToken{}, err
	}
	loaded, loadErr := source.store.LoadAppToken(workContext, key)
	if loadErr != nil && !errors.Is(loadErr, ErrAppTokenNotFound) {
		return AppToken{}, joinContextCause(fmt.Errorf("%w: %w", ErrLoad, loadErr), ctx, workContext)
	}
	if err := checkLease(lease, workContext); err != nil {
		return AppToken{}, err
	}
	if loadErr == nil {
		if err := loaded.Validate(); err != nil {
			return AppToken{}, errors.Join(ErrLoad, err)
		}
		if !loaded.Expired(source.clock.Now(), source.skew) {
			return loaded, nil
		}
	}
	if err := checkLease(lease, workContext); err != nil {
		return AppToken{}, err
	}
	fetched, err := source.fetcher.FetchAppToken(workContext, key)
	if err != nil {
		return AppToken{}, joinContextCause(fmt.Errorf("%w: %w", ErrRefresh, err), ctx, workContext)
	}
	if err := checkLease(lease, workContext); err != nil {
		return AppToken{}, err
	}
	if err := fetched.Validate(); err != nil {
		return AppToken{}, errors.Join(ErrRefresh, err)
	}
	if err := checkLease(lease, workContext); err != nil {
		return AppToken{}, err
	}
	if err := source.store.CommitAppToken(workContext, key, fetched); err != nil {
		return AppToken{}, joinContextCause(fmt.Errorf("%w: %w", ErrCommit, err), ctx, workContext)
	}
	if err := checkLease(lease, workContext); err != nil {
		return AppToken{}, err
	}
	return fetched, nil
}
