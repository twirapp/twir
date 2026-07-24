package channel

import (
	"github.com/google/uuid"
	channelplatform "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type Channel struct {
	ID       uuid.UUID
	ApiKey   *string
	Bindings []channelplatform.ChannelPlatform

	isNil bool
}

func (c Channel) IsNil() bool {
	return c.isNil
}

var Nil = Channel{isNil: true}

func (c Channel) Platforms() []platformentity.Platform {
	platforms := make([]platformentity.Platform, 0, len(c.Bindings))
	for _, binding := range c.Bindings {
		platforms = append(platforms, binding.Platform)
	}

	return platforms
}

func (c Channel) Binding(
	p platformentity.Platform,
) (channelplatform.ChannelPlatform, bool) {
	for _, binding := range c.Bindings {
		if binding.Platform == p {
			return binding, true
		}
	}

	return channelplatform.ChannelPlatform{}, false
}

func (c Channel) BindingByID(
	id uuid.UUID,
) (channelplatform.ChannelPlatform, bool) {
	for _, binding := range c.Bindings {
		if binding.ID == id {
			return binding, true
		}
	}

	return channelplatform.ChannelPlatform{}, false
}

func (c Channel) TwitchBinding() (
	channelplatform.ChannelPlatform,
	channelplatform.TwitchBotConfig,
	bool,
	error,
) {
	binding, found := c.Binding(platformentity.PlatformTwitch)
	if !found {
		return channelplatform.ChannelPlatform{}, channelplatform.TwitchBotConfig{}, false, nil
	}

	config, err := binding.ParseTwitchBotConfig()
	if err != nil {
		return channelplatform.ChannelPlatform{}, channelplatform.TwitchBotConfig{}, false, err
	}

	return binding, config, true, nil
}
