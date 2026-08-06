package postgres

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/spotify_song_request"
)

func TestScanModelToEntity(t *testing.T) {
	now := time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)
	requesterUserID := "user-1"
	requesterDisplayName := "Display"

	model := scanModel{
		ID:                     uuid.MustParse("018f4e3f-1f45-7c4a-8dbb-5e1c5fefc001"),
		ChannelID:              "channel-1",
		TrackID:                "track-1",
		TrackURI:               "spotify:track:track-1",
		Title:                  "Song",
		Artist:                 "Artist",
		Album:                  "Album",
		DurationMs:             123000,
		RequesterUserID:        &requesterUserID,
		RequesterName:          "Requester",
		RequesterDisplayName:   &requesterDisplayName,
		Source:                 "SPOTIFY",
		QueuePosition:          4,
		Status:                 spotify_song_request.StatusPlaying,
		QueuedAt:               now,
		PlayingObservedAt:      &now,
		PlayedObservedAt:       &now,
		CancelledPendingSkipAt: &now,
		SkippedByTwirAt:        &now,
		RemovedOrReconciledAt:  &now,
		UnknownAt:              &now,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	entity := model.toEntity()

	if entity.ID != model.ID || entity.ChannelID != model.ChannelID || entity.TrackID != model.TrackID {
		t.Fatalf("unexpected entity mapping: %+v", entity)
	}
	if entity.RequesterUserID == nil || *entity.RequesterUserID != requesterUserID {
		t.Fatalf("unexpected requester user id: %+v", entity.RequesterUserID)
	}
	if entity.RequesterDisplayName == nil || *entity.RequesterDisplayName != requesterDisplayName {
		t.Fatalf("unexpected requester display name: %+v", entity.RequesterDisplayName)
	}
	if entity.Status != spotify_song_request.StatusPlaying || entity.QueuePosition != 4 || entity.DurationMs != 123000 {
		t.Fatalf("unexpected scalar fields: %+v", entity)
	}
	if entity.PlayingObservedAt == nil || !entity.PlayingObservedAt.Equal(now) {
		t.Fatalf("unexpected playing observed at: %+v", entity.PlayingObservedAt)
	}
}

func TestStatusIsValid(t *testing.T) {
	tests := []struct {
		name   string
		status spotify_song_request.Status
		valid  bool
	}{
		{name: "queued", status: spotify_song_request.StatusQueued, valid: true},
		{name: "playing", status: spotify_song_request.StatusPlaying, valid: true},
		{name: "played", status: spotify_song_request.StatusPlayed, valid: true},
		{name: "cancelled pending skip", status: spotify_song_request.StatusCancelledPendingSkip, valid: true},
		{name: "skipped by twir", status: spotify_song_request.StatusSkippedByTwir, valid: true},
		{name: "removed or reconciled", status: spotify_song_request.StatusRemovedOrReconciled, valid: true},
		{name: "unknown", status: spotify_song_request.StatusUnknown, valid: true},
		{name: "invalid", status: spotify_song_request.Status("invalid"), valid: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.status.IsValid(); got != test.valid {
				t.Fatalf("IsValid() = %v, want %v", got, test.valid)
			}
		})
	}
}
