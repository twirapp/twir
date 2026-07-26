package vkvideo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
)

var (
	ErrInvalidLeaseConfig = errors.New("vk video ownership: invalid lease config")
	ErrLeaseLost          = errors.New("vk video ownership: lease lost")
)

type LeaseConfig struct {
	Expiry        time.Duration
	RenewInterval time.Duration
}

type Ownership struct {
	config  LeaseConfig
	mutexes mutexFactory
	tickers tickerFactory
}

func NewOwnership(redSync *redsync.Redsync, config LeaseConfig) (*Ownership, error) {
	return newOwnership(
		config,
		redSyncMutexFactory{redSync: redSync},
		systemTickerFactory{},
	)
}

func newOwnership(
	config LeaseConfig,
	mutexes mutexFactory,
	tickers tickerFactory,
) (*Ownership, error) {
	if config.Expiry <= 0 {
		return nil, fmt.Errorf("expiry must be positive: %w", ErrInvalidLeaseConfig)
	}
	if config.RenewInterval <= 0 {
		return nil, fmt.Errorf("renew interval must be positive: %w", ErrInvalidLeaseConfig)
	}
	if config.RenewInterval >= config.Expiry {
		return nil, fmt.Errorf("renew interval must be shorter than expiry: %w", ErrInvalidLeaseConfig)
	}

	return &Ownership{
		config:  config,
		mutexes: mutexes,
		tickers: tickers,
	}, nil
}

func (o *Ownership) Acquire(
	ctx context.Context,
	key string,
	closeOwnedResource func(),
) (*Lease, error) {
	mutex := o.mutexes.NewMutex(key, o.config.Expiry)
	if err := mutex.TryLockContext(ctx); err != nil {
		return nil, fmt.Errorf("acquire lease %q: %w", key, err)
	}

	// Bus request handlers are short-lived; lease lifetime is controlled by Release or renewal loss.
	leaseCtx, cancel := context.WithCancelCause(context.WithoutCancel(ctx))
	lease := &Lease{
		ctx:                leaseCtx,
		cancel:             cancel,
		key:                key,
		mutex:              mutex,
		ticker:             o.tickers.NewTicker(o.config.RenewInterval),
		closeOwnedResource: closeOwnedResource,
		stopped:            make(chan struct{}),
	}
	go lease.renew()

	return lease, nil
}

type Lease struct {
	ctx                context.Context
	cancel             context.CancelCauseFunc
	key                string
	mutex              leaseMutex
	ticker             renewalTicker
	closeOwnedResource func()
	stopped            chan struct{}
}

func (l *Lease) Context() context.Context {
	return l.ctx
}

func (l *Lease) Wait() {
	<-l.stopped
}

func (l *Lease) Release(ctx context.Context) error {
	l.cancel(context.Canceled)
	l.Wait()

	released, err := l.mutex.UnlockContext(ctx)
	if !released {
		if err != nil {
			return errors.Join(
				ErrLeaseLost,
				fmt.Errorf("release lease %q: %w", l.key, err),
			)
		}
		return fmt.Errorf("release lease %q: %w", l.key, ErrLeaseLost)
	}
	if err != nil {
		return fmt.Errorf("release lease %q: %w", l.key, err)
	}

	return nil
}

func (l *Lease) renew() {
	defer close(l.stopped)
	defer l.ticker.Stop()
	defer l.closeOwnedResource()

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-l.ticker.Chan():
			renewed, err := l.mutex.ExtendContext(l.ctx)
			if err != nil {
				l.cancel(errors.Join(
					ErrLeaseLost,
					fmt.Errorf("renew lease %q: %w", l.key, err),
				))
				return
			}
			if !renewed {
				l.cancel(fmt.Errorf("renew lease %q: %w", l.key, ErrLeaseLost))
				return
			}
		}
	}
}

type leaseMutex interface {
	TryLockContext(ctx context.Context) error
	ExtendContext(ctx context.Context) (bool, error)
	UnlockContext(ctx context.Context) (bool, error)
}

type mutexFactory interface {
	NewMutex(key string, expiry time.Duration) leaseMutex
}

type redSyncMutexFactory struct {
	redSync *redsync.Redsync
}

func (f redSyncMutexFactory) NewMutex(key string, expiry time.Duration) leaseMutex {
	return f.redSync.NewMutex(key, redsync.WithExpiry(expiry))
}

type renewalTicker interface {
	Chan() <-chan time.Time
	Stop()
}

type tickerFactory interface {
	NewTicker(interval time.Duration) renewalTicker
}

type systemTickerFactory struct{}

func (systemTickerFactory) NewTicker(interval time.Duration) renewalTicker {
	return &systemTicker{Ticker: time.NewTicker(interval)}
}

type systemTicker struct {
	*time.Ticker
}

func (t *systemTicker) Chan() <-chan time.Time {
	return t.C
}
