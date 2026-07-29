package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const renewScript = `if redis.call('get',KEYS[1])==ARGV[1] then return redis.call('pexpire',KEYS[1],ARGV[2]) else return 0 end`
const releaseScript = `if redis.call('get',KEYS[1])==ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end`

type RedisLockerOptions struct {
	Prefix                   string
	TTL, RenewEvery, Timeout time.Duration
}
type RedisLocker struct {
	client  *redis.Client
	options RedisLockerOptions
}

func NewRedisLocker(client *redis.Client, options RedisLockerOptions) (*RedisLocker, error) {
	if client == nil || options.Prefix == "" || options.TTL <= 0 || options.RenewEvery <= 0 || options.RenewEvery >= options.TTL || options.Timeout <= 0 {
		return nil, fmt.Errorf("%w: redis locker options", ErrInvalidOption)
	}
	o := client.Options()
	if !o.ContextTimeoutEnabled || o.ReadTimeout <= 0 || o.WriteTimeout <= 0 || o.PoolTimeout <= 0 {
		return nil, fmt.Errorf("%w: redis client requires context and finite timeouts", ErrInvalidOption)
	}
	return &RedisLocker{client: client, options: options}, nil
}
func (l *RedisLocker) Acquire(ctx context.Context, key CredentialKey) (Lease, error) {
	if err := key.Validate(); err != nil {
		return nil, err
	}
	owner := make([]byte, 32)
	if _, err := rand.Read(owner); err != nil {
		return nil, fmt.Errorf("%w: owner: %w", ErrCoordinator, err)
	}
	value := hex.EncodeToString(owner)
	redisKey := l.options.Prefix + ":" + string(key.Provider) + ":" + string(key.ID)
	for {
		op, cancel := context.WithTimeout(ctx, l.options.Timeout)
		err := l.client.SetArgs(op, redisKey, value, redis.SetArgs{Mode: "NX", TTL: l.options.TTL}).Err()
		cancel()
		if err == nil {
			leaseCtx, cancelLease := context.WithCancelCause(ctx)
			renewCtx, stop := context.WithCancel(leaseCtx)
			lease := &redisLease{client: l.client, key: redisKey, owner: value, ctx: leaseCtx, cancel: cancelLease, renewCtx: renewCtx, stop: stop, options: l.options, done: make(chan struct{})}
			go lease.renew()
			return lease, nil
		}
		if !errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("%w: acquire: %w", ErrCoordinator, err)
		}
		timer := time.NewTimer(l.options.RenewEvery / 4)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}

type redisLease struct {
	client        *redis.Client
	key, owner    string
	ctx, renewCtx context.Context
	cancel        context.CancelCauseFunc
	stop          context.CancelFunc
	options       RedisLockerOptions
	done          chan struct{}
	once          sync.Once
}

func (l *redisLease) Context() context.Context { return l.ctx }
func (l *redisLease) Lost() <-chan struct{}    { return l.ctx.Done() }
func (l *redisLease) renew() {
	defer close(l.done)
	ticker := time.NewTicker(l.options.RenewEvery)
	defer ticker.Stop()
	for {
		select {
		case <-l.renewCtx.Done():
			return
		case <-ticker.C:
			op, cancel := context.WithTimeout(l.renewCtx, l.options.Timeout)
			n, err := l.client.Eval(op, renewScript, []string{l.key}, l.owner, l.options.TTL.Milliseconds()).Int()
			cancel()
			if err != nil || n == 0 {
				if err == nil {
					err = ErrLeaseLost
				}
				l.cancel(errors.Join(ErrLeaseLost, err))
				return
			}
		}
	}
}
func (l *redisLease) Release(ctx context.Context) (err error) {
	l.once.Do(func() {
		l.stop()
		select {
		case <-l.done:
		case <-ctx.Done():
			err = ctx.Err()
			return
		}
		op, cancel := context.WithTimeout(ctx, l.options.Timeout)
		defer cancel()
		n, e := l.client.Eval(op, releaseScript, []string{l.key}, l.owner).Int()
		if e != nil {
			err = fmt.Errorf("%w: release: %w", ErrCoordinator, e)
			return
		}
		if n == 0 {
			l.cancel(ErrLeaseLost)
		}
	})
	return err
}
