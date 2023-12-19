package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type ChannelDuel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	ChannelID string    `gorm:"type:text;"`

	SenderID        null.String `gorm:"type:text;column:sender_id;"`
	SenderModerator bool        `gorm:"type:boolean;column:sender_moderator;"`
	SenderLogin     string      `gorm:"type:text;column:sender_login;"`

	TargetID        null.String `gorm:"type:text;column:target_id;"`
	TargetModerator bool        `gorm:"type:boolean;column:target_moderator;"`
	TargetLogin     string      `gorm:"type:text;column:target_login;"`

	LoserID    null.String `gorm:"type:text;column:loser_id;"`
	CreatedAt  time.Time   `gorm:"type:timestamp;default:now();column:created_at;"`
	FinishedAt null.Time   `gorm:"type:timestamp;column:finished_at;"`

	AvailableUntil time.Time `gorm:"type:timestamp;column:available_until;"`
}

func (c ChannelDuel) TableName() string {
	return "channel_duels"
}
