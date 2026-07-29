package oauth

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	renewLeaseScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
return 0`)
	releaseLeaseScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
end
return 0`)
)

type redisLease struct {
	locker      *RedisLocker
	key         string
	owner       string
	ctx         context.Context
	cancel      context.CancelCauseFunc
	renewCtx    context.Context
	stopRenewal context.CancelFunc
	done        chan struct{}

	mu         sync.Mutex
	lossErr    error
	releaseErr error
	release    sync.Once
}

func newRedisLease(parent context.Context, locker *RedisLocker, key string, owner string) *redisLease {
	leaseContext, cancel := context.WithCancelCause(parent)
	renewContext, stopRenewal := context.WithCancel(leaseContext)
	lease := &redisLease{
		locker: locker, key: key, owner: owner,
		ctx: leaseContext, cancel: cancel,
		renewCtx: renewContext, stopRenewal: stopRenewal,
		done: make(chan struct{}),
	}
	go lease.renew()
	return lease
}

func (lease *redisLease) Context() context.Context { return lease.ctx }
func (lease *redisLease) Lost() <-chan struct{}    { return lease.ctx.Done() }

func (lease *redisLease) Release(ctx context.Context) error {
	if isNil(ctx) {
		return fmt.Errorf("%w: nil release context", ErrInvalidOption)
	}
	lease.release.Do(func() {
		lease.stopRenewal()
		select {
		case <-lease.done:
		case <-ctx.Done():
			lease.releaseErr = lease.markLost(fmt.Errorf("%w: wait for renewal: %w", ErrCoordinator, ctx.Err()))
			return
		}
		if lossErr := lease.loss(); lossErr != nil {
			lease.releaseErr = lossErr
			return
		}
		operationContext, cancel := context.WithTimeout(ctx, lease.locker.options.Timeout)
		defer cancel()
		deleted, err := releaseLeaseScript.Run(operationContext, lease.locker.client, []string{lease.key}, lease.owner).Int64()
		if err != nil {
			lease.releaseErr = lease.markLost(fmt.Errorf("%w: release: %w", ErrCoordinator, err))
			return
		}
		if deleted != 1 {
			lease.releaseErr = lease.markLost(ErrLeaseLost)
			return
		}
		lease.cancel(context.Canceled)
	})
	return lease.releaseErr
}

func (lease *redisLease) renew() {
	ticker := time.NewTicker(lease.locker.options.RenewEvery)
	defer ticker.Stop()
	defer close(lease.done)
	for {
		select {
		case <-lease.renewCtx.Done():
			return
		case <-ticker.C:
			operationContext, cancel := context.WithTimeout(lease.renewCtx, lease.locker.options.Timeout)
			renewed, err := renewLeaseScript.Run(
				operationContext,
				lease.locker.client,
				[]string{lease.key},
				lease.owner,
				lease.locker.options.TTL.Milliseconds(),
			).Int64()
			cancel()
			if err != nil {
				if lease.renewCtx.Err() == nil {
					lease.markLost(fmt.Errorf("%w: renew: %w", ErrCoordinator, err))
				}
				return
			}
			if renewed != 1 {
				lease.markLost(ErrLeaseLost)
				return
			}
		}
	}
}

func (lease *redisLease) markLost(cause error) error {
	lease.mu.Lock()
	newLoss := lease.lossErr == nil
	if lease.lossErr == nil {
		lease.lossErr = errors.Join(ErrLeaseLost, cause)
	}
	lossErr := lease.lossErr
	lease.mu.Unlock()
	if newLoss {
		lease.cancel(lossErr)
	}
	return lossErr
}

func (lease *redisLease) loss() error {
	lease.mu.Lock()
	defer lease.mu.Unlock()
	return lease.lossErr
}
