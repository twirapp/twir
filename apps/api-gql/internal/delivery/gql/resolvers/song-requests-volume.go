package resolvers

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/services/song_requests"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

func getSongRequestSettingsVolume(ctx context.Context, db *gorm.DB, channelID string) int {
	settings := model.ChannelSongRequestsSettings{}
	err := db.WithContext(ctx).
		Select("volume").
		Where(`"channel_id" = ?`, channelID).
		First(&settings).Error
	if err != nil {
		return song_requests.DefaultVolume
	}

	return settings.Volume
}

func saveSongRequestSettingsVolume(ctx context.Context, db *gorm.DB, channelID string, volume int) error {
	return db.WithContext(ctx).Exec(
		`INSERT INTO channels_song_requests_settings (channel_id, enabled, volume)
		VALUES (?, false, ?)
		ON CONFLICT (channel_id) DO UPDATE SET volume = EXCLUDED.volume`,
		channelID,
		volume,
	).Error
}
