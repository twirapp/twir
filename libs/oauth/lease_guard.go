package oauth

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const releaseTimeout = time.Second

func combinedLeaseContext(operationContext context.Context, leaseContext context.Context) (context.Context, func(), error) {
	if isNil(operationContext) || isNil(leaseContext) {
		return nil, func() {}, fmt.Errorf("%w: nil lease context", ErrInvalidOption)
	}
	workContext, cancel := context.WithCancelCause(leaseContext)
	stop := context.AfterFunc(operationContext, func() {
		cause := context.Cause(operationContext)
		if cause == nil {
			cause = operationContext.Err()
		}
		cancel(cause)
	})
	return workContext, func() {
		stop()
		cancel(context.Canceled)
	}, nil
}

func checkLease(lease Lease, workContext context.Context) error {
	if isNil(lease) || isNil(workContext) {
		return fmt.Errorf("%w: invalid lease", ErrCoordinator)
	}
	select {
	case <-workContext.Done():
		cause := context.Cause(workContext)
		if cause == nil {
			cause = workContext.Err()
		}
		if errors.Is(cause, ErrLeaseLost) || !errors.Is(cause, context.Canceled) && !errors.Is(cause, context.DeadlineExceeded) && !errors.Is(cause, ErrClosed) {
			return errors.Join(ErrLeaseLost, cause)
		}
		return cause
	default:
	}
	select {
	case <-lease.Lost():
		return ErrLeaseLost
	default:
		return nil
	}
}

func releaseLease(parent context.Context, lease Lease) error {
	if isNil(parent) || isNil(lease) {
		return fmt.Errorf("%w: invalid release", ErrCoordinator)
	}
	releaseContext, cancel := context.WithTimeout(context.WithoutCancel(parent), releaseTimeout)
	defer cancel()
	if err := lease.Release(releaseContext); err != nil {
		return fmt.Errorf("%w: release: %w", ErrCoordinator, err)
	}
	return nil
}

func joinContextCause(err error, contexts ...context.Context) error {
	joined := err
	for _, ctx := range contexts {
		if isNil(ctx) {
			continue
		}
		if cause := context.Cause(ctx); cause != nil && !errors.Is(joined, cause) {
			joined = errors.Join(joined, cause)
		}
	}
	return joined
}
