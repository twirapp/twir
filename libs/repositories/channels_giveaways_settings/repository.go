package channels_giveaways_settings

import (
	"context"

	"github.com/twirapp/twir/libs/entities/channels_giveaways_settings"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (channels_giveaways_settings.Settings, error)
	Update(ctx context.Context, channelID string, settings channels_giveaways_settings.Settings) (channels_giveaways_settings.Settings, error)
}
