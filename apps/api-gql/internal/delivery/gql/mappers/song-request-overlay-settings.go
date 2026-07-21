package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	entity "github.com/twirapp/twir/libs/entities/song_request_overlay_settings"
)

func SongRequestOverlaySettingsToGQL(
	settings entity.SongRequestOverlaySettings,
) gqlmodel.SongRequestOverlaySettings {
	return gqlmodel.SongRequestOverlaySettings{
		Style:                 gqlmodel.SongRequestOverlayStyle(settings.Style),
		AccentColor:           settings.AccentColor,
		TickerBackgroundColor: settings.TickerBackgroundColor,
		TickerTextColor:       settings.TickerTextColor,
		TickerSpeed:           settings.TickerSpeed,
		HideOnPause:           settings.HideOnPause,
	}
}
