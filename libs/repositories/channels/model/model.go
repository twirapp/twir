package model

import (
	"github.com/google/uuid"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
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

func (c Channel) Platforms() []platformentity.Platform {
	platforms := make([]platformentity.Platform, 0, len(c.Bindings))
	for _, binding := range c.Bindings {
		platforms = append(platforms, binding.Platform)
	}

	return platforms
}

var Nil = Channel{
	isNil: true,
}
