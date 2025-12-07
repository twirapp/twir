package vk_integration

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Entity struct {
	ID          int64
	PublicID    ulid.ULID
	Enabled     bool
	ChannelID   string
	AccessToken string
	UserName    string
	Avatar      string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	isNil bool
}

func (c Entity) IsNil() bool {
	return c.isNil
}

var Nil = Entity{
	isNil: true,
}
