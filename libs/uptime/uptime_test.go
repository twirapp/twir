package uptime

import (
	"testing"
	"time"
)

func TestNewReporter_usesSlotReplica_whenReplicaSet(t *testing.T) {
	// Given
	t.Setenv("REPLICA", "3")

	// When
	reporter, err := NewReporter(nil, ReporterOpts{ServiceName: "eventsub"})

	// Then
	if err != nil {
		t.Fatalf("create reporter: %v", err)
	}
	if got, want := reporter.Instance(), "slot-3"; got != want {
		t.Errorf("instance = %q, want %q", got, want)
	}
}

func TestNewReporter_returnsError_whenServiceNameEmpty(t *testing.T) {
	// Given
	options := ReporterOpts{}

	// When
	_, err := NewReporter(nil, options)

	// Then
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInstanceStatus_Uptime_returnsElapsedDuration(t *testing.T) {
	// Given
	startedAt := time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC)
	status := InstanceStatus{StartedAt: startedAt}
	now := startedAt.Add(2*time.Hour + 15*time.Minute)

	// When
	got := status.Uptime(now)

	// Then
	if want := 2*time.Hour + 15*time.Minute; got != want {
		t.Errorf("uptime = %s, want %s", got, want)
	}
}

func TestInstanceStatus_IsStale_returnsTrue_whenThresholdExceeded(t *testing.T) {
	// Given
	now := time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC)
	status := InstanceStatus{UpdatedAt: now.Add(-61 * time.Second)}

	// When
	got := status.IsStale(now, time.Minute)

	// Then
	if !got {
		t.Error("expected stale status")
	}
}
