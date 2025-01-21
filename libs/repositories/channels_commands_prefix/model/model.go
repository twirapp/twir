package model

import (
	"time"

	"github.com/google/uuid"
)

type ChannelsCommandsPrefix struct {
	ID        uuid.UUID
	ChannelID string
	Prefix    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var Nil = ChannelsCommandsPrefix{}
