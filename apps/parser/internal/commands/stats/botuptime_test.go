package stats

import (
	"testing"
	"time"
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

func TestDisplayInstance_usesSlotNumberWhenReplicaSlot(t *testing.T) {
	// Given
	instance := "slot-3"

	// When
	got := displayInstance(instance)

	// Then
	if want := "#3"; got != want {
		t.Errorf("displayInstance() = %q, want %q", got, want)
	}
}

func TestDisplayInstance_shortensDockerHostnameHash(t *testing.T) {
	// Given
	instance := "a1b2c3d4e5f6"

	// When
	got := displayInstance(instance)

	// Then
	if want := "#4e5f6"; got != want {
		t.Errorf("displayInstance() = %q, want %q", got, want)
	}
}
