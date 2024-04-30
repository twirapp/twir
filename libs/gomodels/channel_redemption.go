package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type ChannelRedemption struct {
	ID           uuid.UUID   `gorm:"type:uuid;column:id"`
	ChannelID    string      `gorm:"type:text;not null;column:channel_id;"`
	UserID       string      `gorm:"type:text;not null;column:user_id"`
	RewardID     uuid.UUID   `gorm:"type:uuid;not null;column:reward_id"`
	RewardTitle  string      `gorm:"type:text;not null;column:reward_title;"`
	RewardPrompt null.String `gorm:"type:text;column:reward_prompt;"`
	RewardCost   int         `gorm:"type:integer;not null;column:reward_cost;"`
	RedeemedAt   time.Time   `gorm:"type:timestamp;not null;column:redeemed_at;"`
}

func (ChannelRedemption) TableName() string {
	return "channel_redemptions_history"
}
