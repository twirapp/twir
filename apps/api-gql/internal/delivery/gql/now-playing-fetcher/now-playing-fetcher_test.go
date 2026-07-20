package now_playing_fetcher

import (
	"testing"
	"time"
)

func intPointer(value int) *int {
	return &value
}

func TestTrackAdvanceProgress(t *testing.T) {
	now := time.Date(2026, time.July, 12, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		track    Track
		expected int
	}{
		{
			name: "adds elapsed cache time",
			track: Track{
				ProgressMs:         intPointer(10000),
				DurationMs:         intPointer(30000),
				ProgressObservedAt: now.Add(-5 * time.Second),
			},
			expected: 15000,
		},
		{
			name: "clamps to duration",
			track: Track{
				ProgressMs:         intPointer(29000),
				DurationMs:         intPointer(30000),
				ProgressObservedAt: now.Add(-5 * time.Second),
			},
			expected: 30000,
		},
		{
			name: "clamps negative progress to zero",
			track: Track{
				ProgressMs:         intPointer(-1000),
				DurationMs:         intPointer(30000),
				ProgressObservedAt: now,
			},
			expected: 0,
		},
		{
			name: "does not subtract progress for future observation",
			track: Track{
				ProgressMs:         intPointer(10000),
				DurationMs:         intPointer(30000),
				ProgressObservedAt: now.Add(5 * time.Second),
			},
			expected: 10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.track.advanceProgress(now)
			if tt.track.ProgressMs == nil || *tt.track.ProgressMs != tt.expected {
				t.Fatalf("expected %d, got %v", tt.expected, tt.track.ProgressMs)
			}
			if !tt.track.ProgressObservedAt.Equal(now) {
				t.Fatal("observation time was not refreshed")
			}
		})
	}
}

func TestTrackAdvanceProgressIgnoresAmbientTrack(t *testing.T) {
	track := Track{}
	track.advanceProgress(time.Now())
	if track.ProgressMs != nil || track.DurationMs != nil {
		t.Fatal("ambient timing must remain nil")
	}
}

func TestTrackBinaryRoundTripPreservesTiming(t *testing.T) {
	observedAt := time.Date(2026, time.July, 12, 12, 0, 0, 0, time.UTC)
	track := Track{
		ProgressMs:         intPointer(10000),
		DurationMs:         intPointer(30000),
		ProgressObservedAt: observedAt,
	}

	data, err := track.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	var cachedTrack Track
	if err := cachedTrack.UnmarshalBinary(data); err != nil {
		t.Fatal(err)
	}

	if cachedTrack.ProgressMs == nil || *cachedTrack.ProgressMs != 10000 {
		t.Fatalf("progress was not preserved: %v", cachedTrack.ProgressMs)
	}
	if cachedTrack.DurationMs == nil || *cachedTrack.DurationMs != 30000 {
		t.Fatalf("duration was not preserved: %v", cachedTrack.DurationMs)
	}
	if !cachedTrack.ProgressObservedAt.Equal(observedAt) {
		t.Fatalf("observation time was not preserved: %v", cachedTrack.ProgressObservedAt)
	}
}
