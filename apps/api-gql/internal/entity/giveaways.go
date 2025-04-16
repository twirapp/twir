package entity

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type ChannelGiveaway struct {
	ID              ulid.ULID
	ChannelID       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	StartedAt       *time.Time
	EndedAt         *time.Time
	IsRunning       bool
	IsStopped       bool
	IsFinished      bool
	Keyword         string
	CreatedByUserID string
	ArchivedAt      *time.Time
	IsArchived      bool
}

var ChannelGiveawayNil = ChannelGiveaway{}
