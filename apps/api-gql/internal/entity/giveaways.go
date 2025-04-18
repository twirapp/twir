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
	Keyword         string
	CreatedByUserID string
	ArchivedAt      *time.Time
	StoppedAt       *time.Time
}

type ChannelGiveawayWinner struct {
	UserID      string
	UserLogin   string
	DisplayName string
}

var ChannelGiveawayNil = ChannelGiveaway{}
