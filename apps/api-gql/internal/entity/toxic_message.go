package entity

import (
	"time"

	"github.com/google/uuid"
)

type ToxicMessage struct {
	ID              uuid.UUID
	ChannelID       *string
	ReplyToToUserID *string
	Text            string
	CreatedAt       time.Time
}
