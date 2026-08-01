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

const expectedRenewScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
return 0`

type fakeRedisCommands struct {
	mu sync.Mutex

	owner       string
	doCalls     [][]any
	evalScripts []string
	evalKeys    [][]string
	evalArgs    [][]any
	evalCtxErrs []error
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
	ctx context.Context,
	script string,
	keys []string,
	args ...any,
) *redis.Cmd {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.evalScripts = append(f.evalScripts, script)
	f.evalKeys = append(f.evalKeys, append([]string(nil), keys...))
	f.evalArgs = append(f.evalArgs, append([]any(nil), args...))
	f.evalCtxErrs = append(f.evalCtxErrs, ctx.Err())

	if script == expectedRenewScript {
		if f.owner == args[0].(string) {
			return redis.NewCmdResult(int64(1), nil)
		}
		return redis.NewCmdResult(int64(0), nil)
	}

	if f.owner == args[0].(string) {
		f.owner = ""
		return redis.NewCmdResult(int64(1), nil)
	}

	return redis.NewCmdResult(int64(0), nil)
}

func (f *fakeRedisCommands) evalContexts() []error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]error(nil), f.evalCtxErrs...)
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

func TestNewRedisRenewsEveryTenSecondsByDefault(t *testing.T) {
	locker := NewRedis(&fakeRedisCommands{})
	if locker.renewInterval != 10*time.Second {
		t.Fatalf("renew interval = %v, want 10s", locker.renewInterval)
	}
	if locker.renewTimeout != 5*time.Second {
		t.Fatalf("renew timeout = %v, want 5s", locker.renewTimeout)
	}
	if locker.leaseWatchdog != 25*time.Second {
		t.Fatalf("lease watchdog = %v, want 25s", locker.leaseWatchdog)
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

type leaseRedisCommands struct {
	mu sync.Mutex

	owner        string
	expiresAt    time.Time
	testLease    time.Duration
	renewCalls   int
	renewNotify  chan struct{}
	loseOnRenew  bool
	renewErr     error
	evalScripts  []string
	renewKeys    [][]string
	renewArgs    [][]any
	releaseCalls int
}

func (f *leaseRedisCommands) Do(_ context.Context, args ...any) *redis.Cmd {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.owner != "" && time.Now().After(f.expiresAt) {
		f.owner = ""
	}
	if f.owner != "" {
		return redis.NewCmdResult(nil, redis.Nil)
	}
	f.owner = args[2].(string)
	f.expiresAt = time.Now().Add(f.testLease)
	return redis.NewCmdResult("OK", nil)
}

func (f *leaseRedisCommands) Eval(
	_ context.Context,
	script string,
	keys []string,
	args ...any,
) *redis.Cmd {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.evalScripts = append(f.evalScripts, script)

	if script == expectedRenewScript {
		f.renewKeys = append(f.renewKeys, append([]string(nil), keys...))
		f.renewArgs = append(f.renewArgs, append([]any(nil), args...))
		f.renewCalls++
		select {
		case f.renewNotify <- struct{}{}:
		default:
		}
		if f.renewErr != nil {
			return redis.NewCmdResult(nil, f.renewErr)
		}
		if f.loseOnRenew {
			f.owner = "replacement-owner"
			return redis.NewCmdResult(int64(0), nil)
		}
		if f.owner == args[0].(string) && time.Now().Before(f.expiresAt) {
			f.expiresAt = time.Now().Add(f.testLease)
			return redis.NewCmdResult(int64(1), nil)
		}
		return redis.NewCmdResult(int64(0), nil)
	}

	f.releaseCalls++
	if f.owner == args[0].(string) {
		f.owner = ""
		return redis.NewCmdResult(int64(1), nil)
	}
	return redis.NewCmdResult(int64(0), nil)
}

func (f *leaseRedisCommands) renewalContract() ([]string, []any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.renewKeys[0]...), append([]any(nil), f.renewArgs[0]...)
}

func (f *leaseRedisCommands) snapshot() (int, []string, int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.renewCalls, append([]string(nil), f.evalScripts...), f.releaseCalls
}

func TestRedisLockerRenewsLeaseWhileCallbackRuns(t *testing.T) {
	commands := &leaseRedisCommands{
		testLease:   40 * time.Millisecond,
		renewNotify: make(chan struct{}, 32),
	}
	locker := newRedis(commands, 5*time.Millisecond)
	firstEntered := make(chan struct{})
	releaseFirst := make(chan struct{})
	firstDone := make(chan error, 1)
	go func() {
		firstDone <- locker.WithLock(context.Background(), "shared", func(context.Context) error {
			close(firstEntered)
			<-releaseFirst
			return nil
		})
	}()
	<-firstEntered

	for range 3 {
		select {
		case <-commands.renewNotify:
		case <-time.After(time.Second):
			t.Fatal("lease was not renewed")
		}
	}

	secondEntered := false
	secondCtx, cancelSecond := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancelSecond()
	err := locker.WithLock(secondCtx, "shared", func(context.Context) error {
		secondEntered = true
		return nil
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("second WithLock() error = %v, want deadline exceeded", err)
	}
	if secondEntered {
		t.Fatal("second owner acquired while first callback lease was being renewed")
	}
	close(releaseFirst)
	if err := <-firstDone; err != nil {
		t.Fatalf("first WithLock() error = %v", err)
	}

	renewCalls, scripts, _ := commands.snapshot()
	if renewCalls < 3 {
		t.Fatalf("renew calls = %d, want at least 3", renewCalls)
	}
	if scripts[0] != expectedRenewScript {
		t.Fatalf("renew script = %q, want exact owner-checked PEXPIRE", scripts[0])
	}
	renewKeys, renewArgs := commands.renewalContract()
	if !reflect.DeepEqual(renewKeys, []string{"shared"}) {
		t.Fatalf("renew keys = %#v, want [shared]", renewKeys)
	}
	if len(renewArgs) != 2 || renewArgs[0] == "" || renewArgs[1] != int64(30000) {
		t.Fatalf("renew args = %#v, want owner and 30000ms lease", renewArgs)
	}
}

func TestRedisLockerCancelsCallbackAndReturnsErrLockLostOnRenewalFailure(t *testing.T) {
	for _, test := range []struct {
		name   string
		result *leaseRedisCommands
	}{
		{
			name: "ownership changed",
			result: &leaseRedisCommands{
				testLease:   time.Second,
				renewNotify: make(chan struct{}, 8),
				loseOnRenew: true,
			},
		},
		{
			name: "redis renewal error",
			result: &leaseRedisCommands{
				testLease:   time.Second,
				renewNotify: make(chan struct{}, 8),
				renewErr:    errors.New("redis unavailable"),
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			locker := newRedis(test.result, 5*time.Millisecond)
			callbackStarted := make(chan struct{})
			callbackExited := make(chan struct{})
			cause := make(chan error, 1)
			var persisted atomic.Bool

			err := locker.WithLock(context.Background(), "shared", func(ctx context.Context) error {
				close(callbackStarted)
				defer close(callbackExited)
				<-ctx.Done()
				cause <- context.Cause(ctx)
				if ctx.Err() == nil {
					persisted.Store(true)
				}
				return ctx.Err()
			})
			if !errors.Is(err, ErrLockLost) {
				t.Fatalf("WithLock() error = %v, want ErrLockLost", err)
			}
			select {
			case <-callbackStarted:
			default:
				t.Fatal("callback did not start")
			}
			select {
			case <-callbackExited:
			default:
				t.Fatal("WithLock() returned before callback exited")
			}
			if got := <-cause; !errors.Is(got, ErrLockLost) {
				t.Fatalf("callback cancellation cause = %v, want ErrLockLost", got)
			}
			if persisted.Load() {
				t.Fatal("callback continued into persistence after ownership loss")
			}

			renewCalls, _, _ := test.result.snapshot()
			time.Sleep(20 * time.Millisecond)
			after, _, _ := test.result.snapshot()
			if after != renewCalls {
				t.Fatalf("renew calls continued after return: before %d after %d", renewCalls, after)
			}
		})
	}
}

func TestRedisLockerReleasesWithDetachedContextAfterCallerCancellation(t *testing.T) {
	commands := &fakeRedisCommands{}
	locker := NewRedis(commands)
	ctx, cancel := context.WithCancel(context.Background())

	err := locker.WithLock(ctx, "shared", func(callbackCtx context.Context) error {
		cancel()
		<-callbackCtx.Done()
		return callbackCtx.Err()
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("WithLock() error = %v, want context canceled", err)
	}

	contexts := commands.evalContexts()
	if len(contexts) != 1 {
		t.Fatalf("Eval calls = %d, want one release", len(contexts))
	}
	if contexts[0] != nil {
		t.Fatalf("release context error = %v, want live detached context", contexts[0])
	}
}

type stallingRenewRedisCommands struct {
	mu sync.Mutex

	owner         string
	renewStarted  chan struct{}
	renewExited   chan struct{}
	startOnce     sync.Once
	exitOnce      sync.Once
	renewCalls    int
	releaseCalls  int
	releaseCtxErr error
}

func newStallingRenewRedisCommands() *stallingRenewRedisCommands {
	return &stallingRenewRedisCommands{
		renewStarted: make(chan struct{}),
		renewExited:  make(chan struct{}),
	}
}

func (f *stallingRenewRedisCommands) Do(_ context.Context, args ...any) *redis.Cmd {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.owner != "" {
		return redis.NewCmdResult(nil, redis.Nil)
	}
	f.owner = args[2].(string)
	return redis.NewCmdResult("OK", nil)
}

func (f *stallingRenewRedisCommands) Eval(
	ctx context.Context,
	script string,
	_ []string,
	args ...any,
) *redis.Cmd {
	if script == expectedRenewScript {
		f.mu.Lock()
		f.renewCalls++
		f.mu.Unlock()
		f.startOnce.Do(func() { close(f.renewStarted) })
		<-ctx.Done()
		f.exitOnce.Do(func() { close(f.renewExited) })
		return redis.NewCmdResult(nil, ctx.Err())
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.releaseCalls++
	f.releaseCtxErr = ctx.Err()
	if f.owner == args[0].(string) {
		f.owner = ""
		return redis.NewCmdResult(int64(1), nil)
	}
	return redis.NewCmdResult(int64(0), nil)
}

func (f *stallingRenewRedisCommands) snapshot() (int, int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.renewCalls, f.releaseCalls, f.releaseCtxErr
}

func TestRedisLockerRenewalDeadlineCancelsCallbackBeforeLeaseExpiry(t *testing.T) {
	commands := newStallingRenewRedisCommands()
	locker := newRedisWithTimings(commands, lockTimings{
		renewInterval: 2 * time.Millisecond,
		renewTimeout:  10 * time.Millisecond,
		leaseWatchdog: 100 * time.Millisecond,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	cause := make(chan error, 1)

	started := time.Now()
	err := locker.WithLock(ctx, "shared", func(ctx context.Context) error {
		<-ctx.Done()
		cause <- context.Cause(ctx)
		return ctx.Err()
	})
	if !errors.Is(err, ErrLockLost) {
		t.Fatalf("WithLock() error = %v, want ErrLockLost", err)
	}
	if elapsed := time.Since(started); elapsed >= 100*time.Millisecond {
		t.Fatalf("WithLock() elapsed = %v, renewal deadline did not return promptly", elapsed)
	}
	if got := <-cause; !errors.Is(got, ErrLockLost) {
		t.Fatalf("callback cause = %v, want ErrLockLost", got)
	}
	select {
	case <-commands.renewExited:
	case <-time.After(time.Second):
		t.Fatal("bounded renewal worker did not exit")
	}
	_, releases, releaseCtxErr := commands.snapshot()
	if releases != 1 || releaseCtxErr != nil {
		t.Fatalf("release = calls %d context error %v, want one detached release", releases, releaseCtxErr)
	}
}

func TestRedisLockerWatchdogCancelsStalledRenewalBeforeLeaseExpiry(t *testing.T) {
	commands := newStallingRenewRedisCommands()
	locker := newRedisWithTimings(commands, lockTimings{
		renewInterval: 2 * time.Millisecond,
		renewTimeout:  120 * time.Millisecond,
		leaseWatchdog: 20 * time.Millisecond,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	cause := make(chan error, 1)

	started := time.Now()
	err := locker.WithLock(ctx, "shared", func(ctx context.Context) error {
		<-ctx.Done()
		cause <- context.Cause(ctx)
		return ctx.Err()
	})
	if !errors.Is(err, ErrLockLost) {
		t.Fatalf("WithLock() error = %v, want ErrLockLost", err)
	}
	if elapsed := time.Since(started); elapsed >= 100*time.Millisecond {
		t.Fatalf("WithLock() elapsed = %v, watchdog did not fire promptly", elapsed)
	}
	if got := <-cause; !errors.Is(got, ErrLockLost) {
		t.Fatalf("callback cause = %v, want ErrLockLost", got)
	}
	select {
	case <-commands.renewExited:
	case <-time.After(time.Second):
		t.Fatal("watchdog cancellation did not stop renewal worker")
	}
}

func TestRedisLockerResetsWatchdogOnlyAfterConfirmedRenewal(t *testing.T) {
	commands := &leaseRedisCommands{
		testLease:   time.Second,
		renewNotify: make(chan struct{}, 32),
	}
	locker := newRedisWithTimings(commands, lockTimings{
		renewInterval: 5 * time.Millisecond,
		renewTimeout:  20 * time.Millisecond,
		leaseWatchdog: 50 * time.Millisecond,
	})

	err := locker.WithLock(context.Background(), "shared", func(context.Context) error {
		for range 12 {
			select {
			case <-commands.renewNotify:
			case <-time.After(time.Second):
				return errors.New("confirmed renewal did not arrive")
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WithLock() error = %v", err)
	}
	renewCalls, _, _ := commands.snapshot()
	if renewCalls < 12 {
		t.Fatalf("renew calls = %d, want at least 12 beyond initial watchdog window", renewCalls)
	}
}

func TestRedisLockerObservesCallbackCompletionDuringStalledRenewal(t *testing.T) {
	commands := newStallingRenewRedisCommands()
	locker := newRedisWithTimings(commands, lockTimings{
		renewInterval: 2 * time.Millisecond,
		renewTimeout:  120 * time.Millisecond,
		leaseWatchdog: 100 * time.Millisecond,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	started := time.Now()
	err := locker.WithLock(ctx, "shared", func(context.Context) error {
		<-commands.renewStarted
		return nil
	})
	if err != nil {
		t.Fatalf("WithLock() error = %v", err)
	}
	if elapsed := time.Since(started); elapsed >= 100*time.Millisecond {
		t.Fatalf("WithLock() elapsed = %v, callback completion was blocked by renewal", elapsed)
	}

	select {
	case <-commands.renewExited:
	case <-time.After(time.Second):
		t.Fatal("canceled renewal worker leaked after callback completion")
	}
	renewCalls, releases, releaseCtxErr := commands.snapshot()
	time.Sleep(15 * time.Millisecond)
	after, _, _ := commands.snapshot()
	if renewCalls != 1 || after != renewCalls {
		t.Fatalf("renew calls = before %d after %d, want one stopped worker", renewCalls, after)
	}
	if releases != 1 || releaseCtxErr != nil {
		t.Fatalf("release = calls %d context error %v, want one detached release", releases, releaseCtxErr)
	}
}
