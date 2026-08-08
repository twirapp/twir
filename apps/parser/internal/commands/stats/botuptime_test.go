package stats

import (
	"testing"
	"time"

	"github.com/twirapp/twir/libs/uptime"
)

func TestFormatUptime_formatsHoursWhenBelowDay(t *testing.T) {
	// Given
	duration := 5*time.Hour + 45*time.Minute

	// When
	got := formatUptime(duration)

	// Then
	if want := "5h"; got != want {
		t.Errorf("formatUptime() = %q, want %q", got, want)
	}
}

func TestInstanceLabels_usesSlotNumberWhenReplicaSlot(t *testing.T) {
	// Given
	instances := []uptime.InstanceStatus{
		{Service: "eventsub", Instance: "slot-3"},
		{Service: "eventsub", Instance: "slot-1"},
	}

	// When
	got := instanceLabels(instances)

	// Then
	if want := "#3"; got["slot-3"] != want {
		t.Errorf("instanceLabels()[slot-3] = %q, want %q", got["slot-3"], want)
	}
	if want := "#1"; got["slot-1"] != want {
		t.Errorf("instanceLabels()[slot-1] = %q, want %q", got["slot-1"], want)
	}
}

func TestInstanceLabels_hidesHostnameForSingleInstance(t *testing.T) {
	// Given
	instances := []uptime.InstanceStatus{
		{Service: "parser", Instance: "rach"},
	}

	// When
	got := instanceLabels(instances)

	// Then
	if want := ""; got["rach"] != want {
		t.Errorf("instanceLabels()[rach] = %q, want %q", got["rach"], want)
	}
}

func TestInstanceLabels_usesOrdinalsForMultipleReplicas(t *testing.T) {
	// Given
	instances := []uptime.InstanceStatus{
		{Service: "parser", Instance: "b2c3d4e5f6a7"},
		{Service: "parser", Instance: "a1b2c3d4e5f6"},
	}

	// When
	got := instanceLabels(instances)

	// Then
	if want := "#1"; got["a1b2c3d4e5f6"] != want {
		t.Errorf("instanceLabels()[a1b2c3d4e5f6] = %q, want %q", got["a1b2c3d4e5f6"], want)
	}
	if want := "#2"; got["b2c3d4e5f6a7"] != want {
		t.Errorf("instanceLabels()[b2c3d4e5f6a7] = %q, want %q", got["b2c3d4e5f6a7"], want)
	}
}

func TestSummarizeServices_ignoresStaleEphemeralInstances(t *testing.T) {
	// Given
	now := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
	statuses := []uptime.InstanceStatus{
		{
			Service:   "parser",
			Instance:  "current-host",
			StartedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Second),
		},
		{
			Service:   "parser",
			Instance:  "old-host",
			StartedAt: now.Add(-2 * time.Hour),
			UpdatedAt: now.Add(-2 * time.Minute),
		},
	}

	// When
	services, down, _ := summarizeServices(statuses, now, "unavailable")

	// Then
	if len(services) != 1 || services[0] != "parser×1 1h" {
		t.Fatalf("services = %v, want [parser×1 1h]", services)
	}
	if len(down) != 0 {
		t.Fatalf("down = %v, want no stale ephemeral instances", down)
	}
}

func TestSummarizeServices_reportsStaleReplicaSlotAsDown(t *testing.T) {
	// Given
	now := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
	statuses := []uptime.InstanceStatus{
		{
			Service:   "eventsub",
			Instance:  "slot-2",
			StartedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-2 * time.Minute),
		},
	}

	// When
	services, down, _ := summarizeServices(statuses, now, "unavailable")

	// Then
	if len(services) != 1 || services[0] != "eventsub×0 unavailable" {
		t.Fatalf("services = %v, want [eventsub×0 unavailable]", services)
	}
	if len(down) != 1 || down[0] != "eventsub#2" {
		t.Fatalf("down = %v, want [eventsub#2]", down)
	}
}
