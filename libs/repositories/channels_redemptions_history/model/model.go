package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
)

type ChannelsRedemptionHistoryItem struct {
	ChannelID    string
	UserID       string
	Platform     platform.Platform
	RewardID     uuid.UUID
	RewardPrompt *string
	RewardTitle  string
	RewardCost   int32
	CreatedAt    time.Time
}
