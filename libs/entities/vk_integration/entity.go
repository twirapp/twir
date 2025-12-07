package vk_integration

import (
	"time"
)

type Entity struct {
	ID          int64
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
