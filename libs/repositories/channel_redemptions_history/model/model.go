package model

import (
	"time"

	"github.com/google/uuid"
)

type ChannelRedemptionHistory struct {
	ID           uuid.UUID
	ChannelID    string
	UserID       string
	RewardID     uuid.UUID
	RewardTitle  string
	RewardPrompt *string
	RewardCost   int
	RedeemedAt   time.Time
}
