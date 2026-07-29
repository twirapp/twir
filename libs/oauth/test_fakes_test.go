package oauth

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type memoryCredentialStore struct {
	mu         sync.Mutex
	credential Credential
	loadErr    error
	commitErr  error
	loadHook   func(context.Context, CredentialKey) error
	commitHook func(context.Context, Credential) error
	loads      int
	commits    int
}

func newMemoryCredentialStore(credential Credential) *memoryCredentialStore {
	return &memoryCredentialStore{credential: cloneTestCredential(credential)}
}

func (s *memoryCredentialStore) Load(ctx context.Context, key CredentialKey) (Credential, error) {
	s.mu.Lock()
	s.loads++
	credential := cloneTestCredential(s.credential)
	err := s.loadErr
	hook := s.loadHook
	s.mu.Unlock()
	if hook != nil {
		if hookErr := hook(ctx, key); hookErr != nil {
			return Credential{}, hookErr
		}
	}
	return credential, err
}

func (s *memoryCredentialStore) Commit(ctx context.Context, credential Credential) error {
	s.mu.Lock()
	s.commits++
	err := s.commitErr
	hook := s.commitHook
	s.mu.Unlock()
	if hook != nil {
		if hookErr := hook(ctx, cloneTestCredential(credential)); hookErr != nil {
			return hookErr
		}
	}
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.credential = cloneTestCredential(credential)
	s.mu.Unlock()
	return nil
}

func (s *memoryCredentialStore) snapshot() (Credential, int, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return cloneTestCredential(s.credential), s.loads, s.commits
}

type recordingRefresher struct {
	mu     sync.Mutex
	result RefreshResult
	err    error
	hook   func(context.Context, Credential) (RefreshResult, error)
	calls  int
}

func (r *recordingRefresher) Refresh(ctx context.Context, credential Credential) (RefreshResult, error) {
	r.mu.Lock()
	r.calls++
	result := cloneTestRefreshResult(r.result)
	err := r.err
	hook := r.hook
	r.mu.Unlock()
	if hook != nil {
		return hook(ctx, cloneTestCredential(credential))
	}
	return result, err
}

func (r *recordingRefresher) callCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.calls
}

type recordingLocker struct {
	lease Lease
	err   error
	calls atomic.Int64
}

func (l *recordingLocker) Acquire(context.Context, CredentialKey) (Lease, error) {
	l.calls.Add(1)
	return l.lease, l.err
}

type controlledLease struct {
	ctx        context.Context
	cancel     context.CancelCauseFunc
	releaseErr error
	released   chan struct{}
	once       sync.Once
}

func newControlledLease(parent context.Context) *controlledLease {
	ctx, cancel := context.WithCancelCause(parent)
	return &controlledLease{ctx: ctx, cancel: cancel, released: make(chan struct{})}
}

func (l *controlledLease) Context() context.Context { return l.ctx }
func (l *controlledLease) Lost() <-chan struct{}    { return l.ctx.Done() }
func (l *controlledLease) Release(context.Context) error {
	l.once.Do(func() {
		l.cancel(context.Canceled)
		close(l.released)
	})
	return l.releaseErr
}
func (l *controlledLease) lose(cause error) { l.cancel(cause) }

type fixedClock struct{ now time.Time }

func (c fixedClock) Now() time.Time { return c.now }

type stagedClock struct {
	mu       sync.Mutex
	now      time.Time
	calls    int
	cancelAt int
	cancel   func()
}

func (c *stagedClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.calls++
	if c.calls == c.cancelAt && c.cancel != nil {
		c.cancel()
	}
	return c.now
}

func cloneTestCredential(credential Credential) Credential {
	if credential.Scopes != nil {
		credential.Scopes = append([]string{}, credential.Scopes...)
	}
	return credential
}

func cloneTestRefreshResult(result RefreshResult) RefreshResult {
	if result.Scopes != nil {
		result.Scopes = append([]string{}, result.Scopes...)
	}
	return result
}
