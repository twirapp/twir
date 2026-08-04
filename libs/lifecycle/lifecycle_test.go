package lifecycle

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestLifecycleRunsHooksInOrder(t *testing.T) {
	t.Parallel()

	lifecycle := New()
	var calls []string
	lifecycle.Append(Hook{
		OnStart: func(context.Context) error {
			calls = append(calls, "start-first")
			return nil
		},
		OnStop: func(context.Context) error {
			calls = append(calls, "stop-first")
			return nil
		},
	})
	lifecycle.Append(Hook{
		OnStart: func(context.Context) error {
			calls = append(calls, "start-second")
			return nil
		},
		OnStop: func(context.Context) error {
			calls = append(calls, "stop-second")
			return nil
		},
	})

	if err := lifecycle.Start(context.Background()); err != nil {
		t.Fatalf("start lifecycle: %v", err)
	}
	if err := lifecycle.Stop(context.Background()); err != nil {
		t.Fatalf("stop lifecycle: %v", err)
	}

	want := []string{"start-first", "start-second", "stop-second", "stop-first"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("unexpected calls: got %v, want %v", calls, want)
	}
}

func TestLifecycleRollsBackStartedHooks(t *testing.T) {
	t.Parallel()

	lifecycle := New()
	var calls []string
	lifecycle.Append(Hook{
		OnStart: func(context.Context) error {
			calls = append(calls, "start-first")
			return nil
		},
		OnStop: func(context.Context) error {
			calls = append(calls, "stop-first")
			return nil
		},
	})
	lifecycle.Append(Hook{
		OnStart: func(context.Context) error {
			calls = append(calls, "start-second")
			return errors.New("boom")
		},
		OnStop: func(context.Context) error {
			calls = append(calls, "stop-second")
			return nil
		},
	})

	if err := lifecycle.Start(context.Background()); err == nil {
		t.Fatal("expected start error")
	}
	if err := lifecycle.Stop(context.Background()); err != nil {
		t.Fatalf("stop lifecycle: %v", err)
	}

	want := []string{"start-first", "start-second", "stop-first"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("unexpected calls: got %v, want %v", calls, want)
	}
}

func TestLifecycleJoinsStopErrors(t *testing.T) {
	t.Parallel()

	firstErr := errors.New("first")
	secondErr := errors.New("second")
	lifecycle := New()
	lifecycle.Append(Hook{OnStop: func(context.Context) error { return firstErr }})
	lifecycle.Append(Hook{OnStop: func(context.Context) error { return secondErr }})

	if err := lifecycle.Start(context.Background()); err != nil {
		t.Fatalf("start lifecycle: %v", err)
	}
	err := lifecycle.Stop(context.Background())
	if !errors.Is(err, firstErr) || !errors.Is(err, secondErr) {
		t.Fatalf("expected joined stop errors, got %v", err)
	}
}
