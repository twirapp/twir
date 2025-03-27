package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type ScheduledVip struct {
	ID        ulid.ULID
	UserID    string
	ChannelID string
	CreatedAt time.Time
	RemoveAt  *time.Time
}

var Nil = ScheduledVip{}
