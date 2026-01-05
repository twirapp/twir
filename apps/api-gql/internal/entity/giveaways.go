package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChannelGiveaway struct {
	ID              uuid.UUID
	ChannelID       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	StartedAt       *time.Time
	Keyword         string
	CreatedByUserID string
	StoppedAt       *time.Time
}

type ChannelGiveawayWinner struct {
	UserID      string
	UserLogin   string
	DisplayName string
}

var ChannelGiveawayNil = ChannelGiveaway{}
