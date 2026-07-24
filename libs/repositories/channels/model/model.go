package model

import (
	"github.com/google/uuid"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
)

type Channel struct {
	ID       uuid.UUID
	ApiKey   *string `db:"api_key"`
	Bindings []channelplatformsmodel.ChannelPlatform

	isNil bool
}

func (c Channel) IsNil() bool {
	return c.isNil
}

var Nil = Channel{
	isNil: true,
}
