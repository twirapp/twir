package model

import (
	"time"

	"github.com/google/uuid"
)

type ChannelsRedemptionHistoryItem struct {
	ChannelID    string
	UserID       string
	RewardID     uuid.UUID
	RewardPrompt *string
	RewardTitle  string
	RewardCost   int32
	CreatedAt    time.Time
}
