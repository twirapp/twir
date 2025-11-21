package entity

import (
	"time"

	"github.com/google/uuid"
)

type TTSUserSettings struct {
	UserID         string
	Rate           int
	Pitch          int
	Volume         int
	Voice          string
	IsChannelOwner bool
}

type BeRightBackOverlay struct {
	ID        uuid.UUID
	ChannelID string
	CreatedAt time.Time
	UpdatedAt time.Time
	Settings  BeRightBackOverlaySettings
}

type BeRightBackOverlaySettings struct {
	Text            string
	Late            BeRightBackOverlayLateSettings
	BackgroundColor string
	FontSize        int32
	FontColor       string
	FontFamily      string
}

type BeRightBackOverlayLateSettings struct {
	Enabled        bool
	Text           string
	DisplayBrbTime bool
}
