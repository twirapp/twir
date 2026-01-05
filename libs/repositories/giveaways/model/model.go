package model

import (
	"time"

	"github.com/google/uuid"
)

type ChannelGiveaway struct {
	ID              uuid.UUID  `db:"id"`
	ChannelID       string     `db:"channel_id"`
	CreatedAt       time.Time  `db:"created_at"`
	Keyword         string     `db:"keyword"`
	UpdatedAt       time.Time  `db:"updated_at"`
	StartedAt       *time.Time `db:"started_at"`
	StoppedAt       *time.Time `db:"stopped_at"`
	CreatedByUserID string     `db:"created_by_user_id"`
}

var ChannelGiveawayNil = ChannelGiveaway{}
