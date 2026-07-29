package oauth

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRedisLockerPausedRedisBoundsRenewerAndRelease(t *testing.T) {
	fixture := newRedisFixture(t)
	lease, err := fixture.locker(t).Acquire(context.Background(), CredentialKey{Provider: "synthetic", ID: "paused"})
	if err != nil {
		t.Fatal(err)
	}
	fixture.docker(t, "pause")
	t.Cleanup(func() {
		fixture.docker(t, "unpause")
		fixture.waitReady(t)
	})
	select {
	case <-lease.Context().Done():
	case <-time.After(time.Second):
		t.Fatal("paused Redis did not stop renewer")
	}
	started := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	err = lease.Release(ctx)
	if !errors.Is(err, ErrLeaseLost) || !errors.Is(err, ErrCoordinator) {
		t.Fatalf("Release error = %v", err)
	}
	if elapsed := time.Since(started); elapsed > 300*time.Millisecond {
		t.Fatalf("Release took %s", elapsed)
	}
	select {
	case <-lease.(*redisLease).done:
	default:
		t.Fatal("renewer remained alive")
	}
}

func TestRedisLockerKilledRedisLosesLeaseAndStopsRenewer(t *testing.T) {
	fixture := newIsolatedRedisFixture(t)
	lease, err := fixture.locker(t).Acquire(context.Background(), CredentialKey{Provider: "synthetic", ID: "killed"})
	if err != nil {
		t.Fatal(err)
	}
	fixture.docker(t, "kill")
	select {
	case <-lease.Context().Done():
	case <-time.After(time.Second):
		t.Fatal("killed Redis did not lose lease")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	err = lease.Release(ctx)
	if !errors.Is(err, ErrLeaseLost) || !errors.Is(err, ErrCoordinator) {
		t.Fatalf("Release error = %v", err)
	}
}
