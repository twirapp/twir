package oauthlock

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const expectedReleaseScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0`

type fakeRedisCommands struct {
	mu sync.Mutex

	owner       string
	doCalls     [][]any
	evalScripts []string
	evalKeys    [][]string
	evalArgs    [][]any
}

func (f *fakeRedisCommands) Do(_ context.Context, args ...any) *redis.Cmd {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.doCalls = append(f.doCalls, append([]any(nil), args...))
	if f.owner != "" {
		return redis.NewCmdResult(nil, redis.Nil)
	}

	f.owner = args[2].(string)
	return redis.NewCmdResult("OK", nil)
}

func (f *fakeRedisCommands) Eval(
	_ context.Context,
	script string,
	keys []string,
	args ...any,
) *redis.Cmd {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.evalScripts = append(f.evalScripts, script)
	f.evalKeys = append(f.evalKeys, append([]string(nil), keys...))
	f.evalArgs = append(f.evalArgs, append([]any(nil), args...))

	if f.owner == args[0].(string) {
		f.owner = ""
		return redis.NewCmdResult(int64(1), nil)
	}

	return redis.NewCmdResult(int64(0), nil)
}

func (f *fakeRedisCommands) snapshot() ([][]any, []string, [][]string, [][]any, string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return append([][]any(nil), f.doCalls...),
		append([]string(nil), f.evalScripts...),
		append([][]string(nil), f.evalKeys...),
		append([][]any(nil), f.evalArgs...),
		f.owner
}

func (f *fakeRedisCommands) replaceOwner(owner string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.owner = owner
}

func TestRedisLockerUsesNXThirtySecondLeaseAndCompareDelete(t *testing.T) {
	commands := &fakeRedisCommands{}
	locker := NewRedis(commands)
	const key = "twir:integration-token-refresh:streamelements:channel-1"

	called := false
	err := locker.WithLock(context.Background(), key, func(context.Context) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("WithLock() error = %v", err)
	}
	if !called {
		t.Fatal("WithLock() did not invoke callback")
	}

	doCalls, scripts, keys, args, owner := commands.snapshot()
	if len(doCalls) != 1 {
		t.Fatalf("SET calls = %d, want 1", len(doCalls))
	}
	if len(doCalls[0]) != 6 {
		t.Fatalf("SET args = %#v, want six arguments", doCalls[0])
	}
	generatedOwner, ok := doCalls[0][2].(string)
	if !ok || generatedOwner == "" {
		t.Fatalf("SET owner = %#v, want generated non-empty UUID", doCalls[0][2])
	}
	if _, err := uuid.Parse(generatedOwner); err != nil {
		t.Fatalf("SET owner = %q, want UUID: %v", generatedOwner, err)
	}
	wantSet := []any{"SET", key, generatedOwner, "NX", "PX", int64(30000)}
	if !reflect.DeepEqual(doCalls[0], wantSet) {
		t.Fatalf("SET args = %#v, want %#v", doCalls[0], wantSet)
	}
	if len(scripts) != 1 || scripts[0] != expectedReleaseScript {
		t.Fatalf("release script = %#v, want exact compare-delete script", scripts)
	}
	if !reflect.DeepEqual(keys[0], []string{key}) {
		t.Fatalf("release keys = %#v, want [%q]", keys[0], key)
	}
	if !reflect.DeepEqual(args[0], []any{generatedOwner}) {
		t.Fatalf("release args = %#v, want owner %q", args[0], generatedOwner)
	}
	if owner != "" {
		t.Fatalf("lock owner after release = %q, want empty", owner)
	}
}

func TestRedisLockerGeneratesANewOwnerForEveryAcquisition(t *testing.T) {
	commands := &fakeRedisCommands{}
	locker := NewRedis(commands)
	for range 2 {
		if err := locker.WithLock(context.Background(), "key", func(context.Context) error {
			return nil
		}); err != nil {
			t.Fatalf("WithLock() error = %v", err)
		}
	}

	doCalls, _, _, _, _ := commands.snapshot()
	if len(doCalls) != 2 {
		t.Fatalf("SET calls = %d, want 2", len(doCalls))
	}
	firstOwner := doCalls[0][2].(string)
	secondOwner := doCalls[1][2].(string)
	if firstOwner == secondOwner {
		t.Fatalf("owners = %q and %q, want unique values", firstOwner, secondOwner)
	}
}

func TestRedisLockerRetriesUntilContextCancellation(t *testing.T) {
	commands := &fakeRedisCommands{owner: "another-process"}
	locker := NewRedis(commands)
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	called := false
	err := locker.WithLock(ctx, "busy", func(context.Context) error {
		called = true
		return nil
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("WithLock() error = %v, want context deadline exceeded", err)
	}
	if called {
		t.Fatal("WithLock() invoked callback without acquiring lock")
	}
	doCalls, scripts, _, _, owner := commands.snapshot()
	if len(doCalls) < 2 {
		t.Fatalf("SET calls = %d, want at least 2 retries", len(doCalls))
	}
	if len(scripts) != 0 {
		t.Fatalf("release scripts = %d, want 0 without ownership", len(scripts))
	}
	if owner != "another-process" {
		t.Fatalf("lock owner = %q, want existing owner preserved", owner)
	}
}

func TestRedisLockerExcludesConcurrentCallbacks(t *testing.T) {
	commands := &fakeRedisCommands{}
	locker := NewRedis(commands)
	firstEntered := make(chan struct{})
	releaseFirst := make(chan struct{})
	secondEntered := make(chan struct{})
	errorsCh := make(chan error, 2)
	var active atomic.Int32
	var overlap atomic.Bool

	go func() {
		errorsCh <- locker.WithLock(context.Background(), "shared", func(context.Context) error {
			if active.Add(1) != 1 {
				overlap.Store(true)
			}
			close(firstEntered)
			<-releaseFirst
			active.Add(-1)
			return nil
		})
	}()
	<-firstEntered

	go func() {
		errorsCh <- locker.WithLock(context.Background(), "shared", func(context.Context) error {
			if active.Add(1) != 1 {
				overlap.Store(true)
			}
			close(secondEntered)
			active.Add(-1)
			return nil
		})
	}()

	select {
	case <-secondEntered:
		t.Fatal("second callback entered while first callback held the lock")
	case <-time.After(50 * time.Millisecond):
	}
	close(releaseFirst)

	select {
	case <-secondEntered:
	case <-time.After(time.Second):
		t.Fatal("second callback did not enter after first released the lock")
	}
	for range 2 {
		if err := <-errorsCh; err != nil {
			t.Fatalf("WithLock() error = %v", err)
		}
	}
	if overlap.Load() {
		t.Fatal("callbacks overlapped")
	}
}

func TestRedisLockerNeverDeletesAnotherOwnersLock(t *testing.T) {
	commands := &fakeRedisCommands{}
	locker := NewRedis(commands)

	err := locker.WithLock(context.Background(), "shared", func(context.Context) error {
		commands.replaceOwner("replacement-owner")
		return nil
	})
	if err != nil {
		t.Fatalf("WithLock() error = %v", err)
	}

	_, scripts, _, _, owner := commands.snapshot()
	if len(scripts) != 1 || scripts[0] != expectedReleaseScript {
		t.Fatalf("release script = %#v, want compare-delete", scripts)
	}
	if owner != "replacement-owner" {
		t.Fatalf("lock owner after release = %q, want replacement owner preserved", owner)
	}
}

func TestRedisLockerReleasesAfterCallbackError(t *testing.T) {
	commands := &fakeRedisCommands{}
	locker := NewRedis(commands)
	wantErr := errors.New("callback failed")

	err := locker.WithLock(context.Background(), "shared", func(context.Context) error {
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("WithLock() error = %v, want callback error", err)
	}

	_, scripts, _, _, owner := commands.snapshot()
	if len(scripts) != 1 {
		t.Fatalf("release scripts = %d, want 1", len(scripts))
	}
	if owner != "" {
		t.Fatalf("lock owner after callback error = %q, want empty", owner)
	}
}
