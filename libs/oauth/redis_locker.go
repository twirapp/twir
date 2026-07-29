package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	credentialLockNamespace = "credential"
	appTokenLockNamespace   = "app"
	acquireRetryFloor       = 10 * time.Millisecond
	acquireRetryCeiling     = 50 * time.Millisecond
)

type RedisLockerOptions struct {
	Prefix     string
	TTL        time.Duration
	RenewEvery time.Duration
	Timeout    time.Duration
	// Hooks must include any source process hooks needed for lease commands because
	// go-redis WithTimeout shares the pool but does not inherit that hook chain.
	Hooks []redis.Hook
}

type RedisLocker struct {
	client  *redis.Client
	options RedisLockerOptions
}

// NewRedisLocker derives bounded lease I/O without mutating or closing client.
// The derived go-redis client shares client's connection pool.
func NewRedisLocker(client *redis.Client, options RedisLockerOptions) (*RedisLocker, error) {
	if client == nil || options.Prefix == "" || options.TTL <= 0 || options.RenewEvery <= 0 ||
		options.RenewEvery > options.TTL/3 || options.Timeout <= 0 || options.Timeout > options.TTL/3 {
		return nil, fmt.Errorf("%w: redis locker options", ErrInvalidOption)
	}
	boundedClient := client.WithTimeout(options.Timeout)
	boundedClient.Options().ContextTimeoutEnabled = true
	for _, hook := range options.Hooks {
		if isNil(hook) {
			return nil, fmt.Errorf("%w: nil Redis hook", ErrInvalidOption)
		}
		boundedClient.AddHook(hook)
	}
	options.Hooks = append([]redis.Hook(nil), options.Hooks...)
	return &RedisLocker{client: boundedClient, options: options}, nil
}

func (locker *RedisLocker) Acquire(ctx context.Context, key CredentialKey) (Lease, error) {
	if isNil(ctx) {
		return nil, fmt.Errorf("%w: nil context", ErrInvalidOption)
	}
	if err := key.Validate(); err != nil {
		return nil, err
	}
	return locker.acquire(ctx, credentialLockNamespace, string(key.Provider), string(key.ID))
}

func (locker *RedisLocker) AcquireAppToken(ctx context.Context, key AppTokenKey) (Lease, error) {
	if isNil(ctx) {
		return nil, fmt.Errorf("%w: nil context", ErrInvalidOption)
	}
	if err := key.Validate(); err != nil {
		return nil, err
	}
	return locker.acquire(ctx, appTokenLockNamespace, string(key.Provider), string(key.ID))
}

func (locker *RedisLocker) acquire(ctx context.Context, namespace string, provider string, id string) (Lease, error) {
	owner, err := newLeaseOwner()
	if err != nil {
		return nil, fmt.Errorf("%w: owner: %w", ErrCoordinator, err)
	}
	key := locker.redisKey(namespace, provider, id)
	for {
		if cause := context.Cause(ctx); cause != nil {
			return nil, cause
		}
		operationContext, cancel := context.WithTimeout(ctx, locker.options.Timeout)
		acquired, operationErr := locker.client.SetNX(operationContext, key, owner, locker.options.TTL).Result()
		cancel()
		if operationErr != nil {
			return nil, fmt.Errorf("%w: acquire: %w", ErrCoordinator, operationErr)
		}
		if acquired {
			return newRedisLease(ctx, locker, key, owner), nil
		}
		delay, err := acquireDelay()
		if err != nil {
			return nil, fmt.Errorf("%w: jitter: %w", ErrCoordinator, err)
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return nil, context.Cause(ctx)
		case <-timer.C:
		}
	}
}

func (locker *RedisLocker) redisKey(namespace string, provider string, id string) string {
	encode := base64.RawURLEncoding.EncodeToString
	return locker.options.Prefix + ":" + namespace + ":" + encode([]byte(provider)) + ":" + encode([]byte(id))
}

func newLeaseOwner() (string, error) {
	owner := make([]byte, 32)
	if _, err := rand.Read(owner); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(owner), nil
}

func acquireDelay() (time.Duration, error) {
	span := int64(acquireRetryCeiling - acquireRetryFloor)
	random, err := rand.Int(rand.Reader, big.NewInt(span+1))
	if err != nil {
		return 0, err
	}
	return acquireRetryFloor + time.Duration(random.Int64()), nil
}
