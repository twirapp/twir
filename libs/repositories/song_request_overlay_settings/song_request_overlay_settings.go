package song_request_overlay_settings

import (
	"context"

	entity "github.com/twirapp/twir/libs/entities/song_request_overlay_settings"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (entity.SongRequestOverlaySettings, error)
	Upsert(ctx context.Context, input UpsertInput) (entity.SongRequestOverlaySettings, error)
}

type UpsertInput struct {
	ChannelID             string
	Style                 entity.Style
	AccentColor           string
	TickerBackgroundColor string
	TickerTextColor       string
	TickerSpeed           int
	HideOnPause           bool
}
