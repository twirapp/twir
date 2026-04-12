package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
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
	RolesIDs         []uuid.UUID
	Platforms        []platform.Platform
}

var Nil = Keyword{}
