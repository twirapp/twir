package pgx

import (
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	"github.com/twirapp/twir/libs/repositories/channels/model"
)

func mapBindingToEntity(m channelplatformsmodel.ChannelPlatform) channelplatformentity.ChannelPlatform {
	if m.IsNil() {
		return channelplatformentity.Nil
	}

	return channelplatformentity.ChannelPlatform{
		ID:                m.ID,
		ChannelID:         m.ChannelID,
		Platform:          m.Platform,
		UserID:            m.UserID,
		PlatformChannelID: m.PlatformChannelID,
		Enabled:           m.Enabled,
		BotUserID:         m.BotUserID,
		BotConfig:         m.BotConfig,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func mapChannelToEntity(m model.Channel) channelentity.Channel {
	if m.IsNil() {
		return channelentity.Nil
	}

	bindings := make([]channelplatformentity.ChannelPlatform, 0, len(m.Bindings))
	for _, binding := range m.Bindings {
		bindings = append(bindings, mapBindingToEntity(binding))
	}

	return channelentity.Channel{
		ID:       m.ID,
		ApiKey:   m.ApiKey,
		Bindings: bindings,
	}
}
