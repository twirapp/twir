package model

import (
	"time"

	"github.com/google/uuid"
)

type BeRightBackOverlay struct {
	ID        uuid.UUID
	ChannelID string
	CreatedAt time.Time
	UpdatedAt time.Time
	Settings  *BeRightBackOverlaySettings
}

type BeRightBackOverlaySettings struct {
	Text            string                         `json:"text"`
	Late            BeRightBackOverlayLateSettings `json:"late"`
	BackgroundColor string                         `json:"background_color"`
	FontSize        int32                          `json:"font_size"`
	FontColor       string                         `json:"font_color"`
	FontFamily      string                         `json:"font_family"`
}

type BeRightBackOverlayLateSettings struct {
	Enabled        bool   `json:"enabled"`
	Text           string `json:"text"`
	DisplayBrbTime bool   `json:"display_brb_time"`
}

var Nil = BeRightBackOverlay{}
