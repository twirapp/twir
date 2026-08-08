package spotify_song_request

import (
	"database/sql/driver"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type SpotifySongRequest struct {
	ID                     uuid.UUID
	ChannelID              string
	TrackID                string
	TrackURI               string
	Title                  string
	Artist                 string
	Album                  string
	DurationMs             int
	RequesterUserID        *string
	RequesterName          string
	RequesterDisplayName   *string
	Source                 string
	QueuePosition          int
	Status                 Status
	QueuedAt               time.Time
	PlayingObservedAt      *time.Time
	PlayedObservedAt       *time.Time
	CancelledPendingSkipAt *time.Time
	SkippedByTwirAt        *time.Time
	RemovedOrReconciledAt  *time.Time
	UnknownAt              *time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time

	isNil bool
}

func (c SpotifySongRequest) IsNil() bool {
	return c.isNil
}

var Nil = SpotifySongRequest{isNil: true}

type Status string

const (
	StatusQueued               Status = "queued"
	StatusPlaying              Status = "playing"
	StatusPlayed               Status = "played"
	StatusCancelledPendingSkip Status = "cancelled_pending_skip"
	StatusSkippedByTwir        Status = "skipped_by_twir"
	StatusRemovedOrReconciled  Status = "removed_or_reconciled"
	StatusUnknown              Status = "unknown"
)

var allStatuses = []Status{
	StatusQueued,
	StatusPlaying,
	StatusPlayed,
	StatusCancelledPendingSkip,
	StatusSkippedByTwir,
	StatusRemovedOrReconciled,
	StatusUnknown,
}

func (s Status) IsValid() bool {
	return slices.Contains(allStatuses, s)
}

func (s Status) String() string {
	return string(s)
}

func (s *Status) Scan(src any) error {
	switch v := src.(type) {
	case string:
		*s = Status(v)
	case []byte:
		*s = Status(v)
	case nil:
		*s = ""
	default:
		return fmt.Errorf("spotify song request status: cannot scan type %T into Status", src)
	}

	return nil
}

func (s Status) Value() (driver.Value, error) {
	return string(s), nil
}

var ErrInvalidStatus = fmt.Errorf("invalid spotify song request status")
