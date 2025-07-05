package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChannelRedemptionHistoryItem struct {
	ChannelID    string
	UserID       string
	RewardID     uuid.UUID
	RewardPrompt *string
	RewardTitle  string
	RewardCost   int32
	RedeemedAt   time.Time
}
