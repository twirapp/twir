package model

import (
	"time"

	"github.com/google/uuid"
)

type ChannelsCommandsUsages struct {
	UserID    string
	ChannelID string
	CommandID uuid.UUID
	CreatedAt time.Time
}
