package oauthlock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	lockTTL       = 30 * time.Second
	renewInterval = 10 * time.Second
	retryDelay    = 25 * time.Millisecond
	releaseLimit  = 5 * time.Second
)

const releaseScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0`

const renewScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
return 0`

var ErrLockLost = errors.New("oauth refresh lock lost")

type Locker interface {
	WithLock(ctx context.Context, key string, fn func(context.Context) error) error
}

type RedisCommands interface {
	Do(ctx context.Context, args ...any) *redis.Cmd
	Eval(ctx context.Context, script string, keys []string, args ...any) *redis.Cmd
}

type RedisLocker struct {
	commands      RedisCommands
	renewInterval time.Duration
}

func NewRedis(commands RedisCommands) *RedisLocker {
	return newRedis(commands, renewInterval)
}

func newRedis(commands RedisCommands, interval time.Duration) *RedisLocker {
	return &RedisLocker{
		commands:      commands,
		renewInterval: interval,
	}
}

func (l *RedisLocker) WithLock(
	ctx context.Context,
	key string,
	fn func(context.Context) error,
) error {
	owner := uuid.NewString()

	for {
		_, err := l.commands.Do(
			ctx,
			"SET",
			key,
			owner,
			"NX",
			"PX",
			lockTTL.Milliseconds(),
		).Text()
		switch {
		case err == nil:
			return l.runCallback(ctx, key, owner, fn)
		case !errors.Is(err, redis.Nil):
			return fmt.Errorf("acquire OAuth refresh lock: %w", err)
		}

		timer := time.NewTimer(retryDelay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return ctx.Err()
		case <-timer.C:
		}
	}
}

func (l *RedisLocker) runCallback(
	ctx context.Context,
	key, owner string,
	fn func(context.Context) error,
) error {
	callbackCtx, cancelCallback := context.WithCancelCause(ctx)
	defer cancelCallback(nil)

	callbackDone := make(chan error, 1)
	go func() {
		callbackDone <- fn(callbackCtx)
	}()

	ticker := time.NewTicker(l.renewInterval)
	defer ticker.Stop()

	finish := func(primaryErr, callbackErr error) error {
		return errors.Join(primaryErr, callbackErr, l.release(ctx, key, owner))
	}

	for {
		select {
		case callbackErr := <-callbackDone:
			return finish(nil, callbackErr)
		default:
		}

		select {
		case callbackErr := <-callbackDone:
			return finish(nil, callbackErr)
		case <-ctx.Done():
			cancelCallback(context.Cause(ctx))
			return finish(ctx.Err(), <-callbackDone)
		case <-ticker.C:
			if err := ctx.Err(); err != nil {
				cancelCallback(context.Cause(ctx))
				return finish(err, <-callbackDone)
			}

			renewed, err := l.renew(callbackCtx, key, owner)
			if err == nil && renewed {
				continue
			}

			cancelCallback(ErrLockLost)
			callbackErr := <-callbackDone
			if err != nil {
				err = fmt.Errorf("renew OAuth refresh lock: %w", err)
			}
			return finish(errors.Join(ErrLockLost, err), callbackErr)
		}
	}
}

func (l *RedisLocker) renew(ctx context.Context, key, owner string) (bool, error) {
	result, err := l.commands.Eval(
		ctx,
		renewScript,
		[]string{key},
		owner,
		lockTTL.Milliseconds(),
	).Int64()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (l *RedisLocker) release(ctx context.Context, key, owner string) error {
	releaseCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), releaseLimit)
	defer cancel()

	if err := l.commands.Eval(releaseCtx, releaseScript, []string{key}, owner).Err(); err != nil {
		return fmt.Errorf("release OAuth refresh lock: %w", err)
	}

	return nil
}
