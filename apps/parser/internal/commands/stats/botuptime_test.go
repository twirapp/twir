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
