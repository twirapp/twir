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
	lockTTL      = 30 * time.Second
	retryDelay   = 25 * time.Millisecond
	releaseLimit = 5 * time.Second
)

const releaseScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0`

type Locker interface {
	WithLock(ctx context.Context, key string, fn func(context.Context) error) error
}

type RedisCommands interface {
	Do(ctx context.Context, args ...any) *redis.Cmd
	Eval(ctx context.Context, script string, keys []string, args ...any) *redis.Cmd
}

type RedisLocker struct {
	commands RedisCommands
}

func NewRedis(commands RedisCommands) *RedisLocker {
	return &RedisLocker{commands: commands}
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
			callbackErr := fn(ctx)
			releaseErr := l.release(ctx, key, owner)
			return errors.Join(callbackErr, releaseErr)
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

func (l *RedisLocker) release(ctx context.Context, key, owner string) error {
	releaseCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), releaseLimit)
	defer cancel()

	if err := l.commands.Eval(releaseCtx, releaseScript, []string{key}, owner).Err(); err != nil {
		return fmt.Errorf("release OAuth refresh lock: %w", err)
	}

	return nil
}
