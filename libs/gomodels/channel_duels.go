package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type ChannelDuel struct {
	ID        uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4()"`
	ChannelID string      `gorm:"type:text;"`
	SenderID  null.String `gorm:"type:text;"`
	TargetID  null.String `gorm:"type:text;"`
	LoserID   null.String `gorm:"type:text;"`
	CreatedAt time.Time   `gorm:"type:timestamp;default:now()"`
}

func (c ChannelDuel) TableName() string {
	return "channel_duels"
}
