package oauth

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestRedisLockerDerivesBoundedContextAwareClient(t *testing.T) {
	fixture := newRedisFixture(t)
	sourceOptions := *fixture.source.Options()
	locker := fixture.locker(t)

	if locker.client == fixture.source {
		t.Fatal("locker did not derive a client")
	}
	options := locker.client.Options()
	if !options.ContextTimeoutEnabled || options.ReadTimeout <= 0 || options.WriteTimeout <= 0 {
		t.Fatalf("derived options are not bounded: %+v", options)
	}
	if fixture.source.Options().ContextTimeoutEnabled != sourceOptions.ContextTimeoutEnabled ||
		fixture.source.Options().ReadTimeout != sourceOptions.ReadTimeout ||
		fixture.source.Options().WriteTimeout != sourceOptions.WriteTimeout {
		t.Fatal("source client options were changed")
	}
}

func TestRedisLockerInstallsExplicitHooksOnDerivedClient(t *testing.T) {
	fixture := newRedisFixture(t)
	hook := &countingRedisHook{}
	locker, err := NewRedisLocker(fixture.source, RedisLockerOptions{
		Prefix: "oauth-acceptance", TTL: 600 * time.Millisecond,
		RenewEvery: 150 * time.Millisecond, Timeout: 100 * time.Millisecond,
		Hooks: []redis.Hook{hook},
	})
	if err != nil {
		t.Fatal(err)
	}
	lease, err := locker.Acquire(context.Background(), CredentialKey{Provider: "synthetic", ID: "hook"})
	if err != nil {
		t.Fatal(err)
	}
	if err := lease.Release(context.Background()); err != nil {
		t.Fatal(err)
	}
	if hook.calls.Load() < 2 {
		t.Fatalf("hook calls = %d", hook.calls.Load())
	}
}

func TestRedisLockerEscapesCredentialKeySegments(t *testing.T) {
	fixture := newRedisFixture(t)
	locker := fixture.locker(t)
	first, err := locker.Acquire(context.Background(), CredentialKey{Provider: "a:b", ID: "c"})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = first.Release(context.Background()) }()
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	second, err := locker.Acquire(ctx, CredentialKey{Provider: "a", ID: "b:c"})
	if err != nil {
		t.Fatalf("escaped keys collided: %v", err)
	}
	if err := second.Release(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestRedisLockerWrongOwnerReleaseRetainsKeyAndLosesLease(t *testing.T) {
	fixture := newRedisFixture(t)
	locker := fixture.locker(t)
	lease, err := locker.Acquire(context.Background(), CredentialKey{Provider: "synthetic", ID: "wrong-owner"})
	if err != nil {
		t.Fatal(err)
	}
	keys := fixture.keys(t)
	if len(keys) != 1 {
		t.Fatalf("keys = %#v", keys)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := fixture.admin.Set(ctx, keys[0], "replacement-owner", time.Second).Err(); err != nil {
		t.Fatal(err)
	}
	if err := lease.Release(ctx); !errors.Is(err, ErrLeaseLost) {
		t.Fatalf("Release error = %v", err)
	}
	value, err := fixture.admin.Get(ctx, keys[0]).Result()
	if err != nil || value != "replacement-owner" {
		t.Fatalf("replacement owner = %q, error = %v", value, err)
	}
	select {
	case <-lease.Context().Done():
	default:
		t.Fatal("wrong-owner release left lease context alive")
	}
}

func TestRedisLockerSuccessfulReleaseIsConcurrentAndClosesContext(t *testing.T) {
	fixture := newRedisFixture(t)
	lease, err := fixture.locker(t).Acquire(context.Background(), CredentialKey{Provider: "synthetic", ID: "concurrent"})
	if err != nil {
		t.Fatal(err)
	}
	var group sync.WaitGroup
	errorsChannel := make(chan error, 16)
	for range 16 {
		group.Add(1)
		go func() {
			defer group.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			errorsChannel <- lease.Release(ctx)
		}()
	}
	group.Wait()
	close(errorsChannel)
	for err := range errorsChannel {
		if err != nil {
			t.Fatalf("Release error = %v", err)
		}
	}
	select {
	case <-lease.Context().Done():
	default:
		t.Fatal("successful release left lease context alive")
	}
	if keys := fixture.keys(t); len(keys) != 0 {
		t.Fatalf("released keys = %#v", keys)
	}
}

func TestRedisLockerRenewsTTLAndUsesUniqueCryptographicOwners(t *testing.T) {
	fixture := newRedisFixture(t)
	locker := fixture.locker(t)
	key := CredentialKey{Provider: "synthetic", ID: "renewed"}
	first, err := locker.Acquire(context.Background(), key)
	if err != nil {
		t.Fatal(err)
	}
	keys := fixture.keys(t)
	if len(keys) != 1 {
		t.Fatalf("keys = %#v", keys)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	firstOwner, err := fixture.admin.Get(ctx, keys[0]).Result()
	if err != nil || len(firstOwner) < 43 {
		t.Fatalf("owner generation failed: length = %d, error = %v", len(firstOwner), err)
	}
	timer := time.NewTimer(750 * time.Millisecond)
	defer timer.Stop()
	<-timer.C
	if exists, err := fixture.admin.Exists(ctx, keys[0]).Result(); err != nil || exists != 1 {
		t.Fatalf("renewed key exists = %d, error = %v", exists, err)
	}
	if err := first.Release(ctx); err != nil {
		t.Fatal(err)
	}
	second, err := locker.Acquire(context.Background(), key)
	if err != nil {
		t.Fatal(err)
	}
	secondOwner, err := fixture.admin.Get(ctx, keys[0]).Result()
	if err != nil {
		t.Fatal(err)
	}
	if firstOwner == secondOwner {
		t.Fatal("lease owner was reused")
	}
	if err := second.Release(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestRedisLockerNilContextsReturnErrors(t *testing.T) {
	fixture := newRedisFixture(t)
	locker := fixture.locker(t)
	var nilContext *panicContext
	if _, err := locker.Acquire(nilContext, CredentialKey{Provider: "synthetic", ID: "nil-context"}); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("Acquire error = %v", err)
	}
	lease, err := locker.Acquire(context.Background(), CredentialKey{Provider: "synthetic", ID: "nil-release"})
	if err != nil {
		t.Fatal(err)
	}
	if err := lease.Release(nilContext); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("Release error = %v", err)
	}
	if err := lease.Release(context.Background()); err != nil {
		t.Fatal(err)
	}
}

type countingRedisHook struct{ calls atomic.Int64 }

func (*countingRedisHook) DialHook(next redis.DialHook) redis.DialHook { return next }
func (hook *countingRedisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, command redis.Cmder) error {
		hook.calls.Add(1)
		return next(ctx, command)
	}
}
func (*countingRedisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}
