package model

import (
	"time"

	"github.com/google/uuid"
)

type ToxicMessage struct {
	ID              uuid.UUID
	ChannelID       *string
	ReplyToToUserID *string `db:"reply_to_user_id"`
	Text            string
	CreatedAt       time.Time
}
