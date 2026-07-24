package pgx

import (
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/repositories/channel_platforms/model"
)

func mapBindingToEntity(m model.ChannelPlatform) channelplatformentity.ChannelPlatform {
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

func mapBindingsToEntities(models []model.ChannelPlatform) []channelplatformentity.ChannelPlatform {
	entities := make([]channelplatformentity.ChannelPlatform, 0, len(models))
	for _, m := range models {
		entities = append(entities, mapBindingToEntity(m))
	}

	return entities
}
