package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const defaultTimeout = 15 * time.Second

// Hook defines bounded work performed while an application starts or stops.
type Hook struct {
	OnStart func(context.Context) error
	OnStop  func(context.Context) error
}

// Lifecycle owns application startup and shutdown hooks.
type Lifecycle struct {
	mu      sync.Mutex
	hooks   []Hook
	started int
	running bool
	stopped bool
}

func New() *Lifecycle {
	return &Lifecycle{}
}

// Append registers a hook. Applications must finish registration before Start.
func (l *Lifecycle) Append(hook Hook) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.running || l.stopped {
		panic("lifecycle: append after start")
	}

	l.hooks = append(l.hooks, hook)
}

// Start runs hooks in registration order.
func (l *Lifecycle) Start(ctx context.Context) error {
	l.mu.Lock()
	if l.running || l.stopped {
		l.mu.Unlock()
		return errors.New("lifecycle: already started")
	}
	l.running = true
	hooks := append([]Hook(nil), l.hooks...)
	l.mu.Unlock()

	for index, hook := range hooks {
		if hook.OnStart != nil {
			if err := hook.OnStart(ctx); err != nil {
				l.mu.Lock()
				l.started = index
				l.mu.Unlock()
				return fmt.Errorf("start hook %d: %w", index, err)
			}
		}

		l.mu.Lock()
		l.started = index + 1
		l.mu.Unlock()
	}

	return nil
}

// Stop runs hooks in reverse registration order.
func (l *Lifecycle) Stop(ctx context.Context) error {
	l.mu.Lock()
	if l.stopped {
		l.mu.Unlock()
		return nil
	}
	l.stopped = true
	started := l.started
	hooks := append([]Hook(nil), l.hooks[:started]...)
	l.mu.Unlock()

	var errs []error
	for index := len(hooks) - 1; index >= 0; index-- {
		if hooks[index].OnStop == nil {
			continue
		}

		if err := hooks[index].OnStop(ctx); err != nil {
			errs = append(errs, fmt.Errorf("stop hook %d: %w", index, err))
		}
	}

	return errors.Join(errs...)
}

// Run starts the application, waits for SIGINT or SIGTERM, and stops it.
func (l *Lifecycle) Run() error {
	startCtx, cancelStart := context.WithTimeout(context.Background(), defaultTimeout)
	err := l.Start(startCtx)
	cancelStart()
	if err != nil {
		stopCtx, cancelStop := context.WithTimeout(context.Background(), defaultTimeout)
		stopErr := l.Stop(stopCtx)
		cancelStop()
		return errors.Join(err, stopErr)
	}

	signalCtx, cancelSignal := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	<-signalCtx.Done()
	cancelSignal()

	stopCtx, cancelStop := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancelStop()

	return l.Stop(stopCtx)
}
