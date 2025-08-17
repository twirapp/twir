package model

import (
	"time"

	"github.com/google/uuid"
)

type UserStat struct {
	ID                uuid.UUID
	UserID            string
	ChannelID         string
	Messages          int32
	Watched           int64
	UsedChannelPoints int64
	IsMod             bool
	IsVip             bool
	IsSubscriber      bool
	Reputation        int64
	Emotes            int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
