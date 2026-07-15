package song_request_overlay_settings

import (
	"time"

	"github.com/google/uuid"
)

type Style string

const (
	StyleCinema   Style = "CINEMA"
	StyleCompact  Style = "COMPACT"
	StyleTicker   Style = "TICKER"
	StyleStudio   Style = "STUDIO"
	StylePortrait Style = "PORTRAIT"
	StylePill     Style = "PILL"

	DefaultAccentColor           = "#8B5CF6"
	DefaultTickerBackgroundColor = "#111827E6"
	DefaultTickerTextColor       = "#FFFFFF"
	DefaultTickerSpeed           = 35
)

func (s Style) IsValid() bool {
	switch s {
	case StyleCinema, StyleCompact, StyleTicker, StyleStudio, StylePortrait, StylePill:
		return true
	default:
		return false
	}
}

type SongRequestOverlaySettings struct {
	ID                    uuid.UUID
	ChannelID             string
	Style                 Style
	AccentColor           string
	TickerBackgroundColor string
	TickerTextColor       string
	TickerSpeed           int
	HideOnPause           bool
	CreatedAt             time.Time
	UpdatedAt             time.Time

	isNil bool
}

func (s SongRequestOverlaySettings) IsNil() bool {
	return s.isNil
}

func Default(channelID string) SongRequestOverlaySettings {
	return SongRequestOverlaySettings{
		ChannelID:             channelID,
		Style:                 StyleCinema,
		AccentColor:           DefaultAccentColor,
		TickerBackgroundColor: DefaultTickerBackgroundColor,
		TickerTextColor:       DefaultTickerTextColor,
		TickerSpeed:           DefaultTickerSpeed,
		HideOnPause:           true,
	}
}

var Nil = SongRequestOverlaySettings{isNil: true}
