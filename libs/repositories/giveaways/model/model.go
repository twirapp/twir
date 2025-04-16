package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type ChannelGiveaway struct {
	ID              ulid.ULID  `db:"id"`
	ChannelID       string     `db:"channel_id"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	StartedAt       *time.Time `db:"started_at"`
	ArchivedAt      *time.Time
	EndedAt         *time.Time `db:"ended_at"`
	IsRunning       bool       `db:"is_running"`
	IsStopped       bool       `db:"is_stopped"`
	IsFinished      bool       `db:"is_finished"`
	IsArchived      bool       `db:"is_archived"`
	Keyword         string     `db:"keyword"`
	CreatedByUserID string     `db:"created_by_user_id"`
}

var ChannelGiveawayNil = ChannelGiveaway{}
