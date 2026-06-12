package model

import (
	"time"

	"github.com/google/uuid"
)

type UserStat struct {
	ID                uuid.UUID `db:"id"`
	UserID            uuid.UUID `db:"user_id"`
	ChannelID         uuid.UUID `db:"channel_id"`
	Messages          int32     `db:"messages"`
	Watched           int64     `db:"watched"`
	UsedChannelPoints int64     `db:"usedChannelPoints"`
	IsMod             bool      `db:"is_mod"`
	IsVip             bool      `db:"is_vip"`
	IsSubscriber      bool      `db:"is_subscriber"`
	Reputation        int64     `db:"reputation"`
	Emotes            int       `db:"emotes"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}
