package pgx

import (
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
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
