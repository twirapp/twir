package song_requests_settings

import (
	"context"

	"github.com/twirapp/twir/libs/entities/song_requests_settings"
)

type Repository interface {
	GetByChannelID(
		ctx context.Context,
		channelID string,
	) (song_requests_settings.Settings, error)
	Upsert(
		ctx context.Context,
		settings song_requests_settings.Settings,
	) (song_requests_settings.Settings, error)
	SetVolume(ctx context.Context, channelID string, volume int) error
}
