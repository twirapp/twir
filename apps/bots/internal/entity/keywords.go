package entity

import (
	"time"

	"github.com/google/uuid"
)

type Keyword struct {
	ID               uuid.UUID
	ChannelID        string
	Text             string
	Response         string
	Enabled          bool
	Cooldown         int
	CooldownExpireAt *time.Time
	IsReply          bool
	IsRegular        bool
	Usages           int
}
