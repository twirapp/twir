package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChannelOverlayType string

const (
	ChannelOverlayTypeHTML  ChannelOverlayType = "HTML"
	ChannelOverlayTypeIMAGE ChannelOverlayType = "IMAGE"
)

type ChannelOverlayLayerSettings struct {
	HtmlOverlayHTML                    string
	HtmlOverlayCSS                     string
	HtmlOverlayJS                      string
	HtmlOverlayDataPollSecondsInterval int
	ImageUrl                           string
}

type ChannelOverlayLayer struct {
	ID                      uuid.UUID
	Type                    ChannelOverlayType
	Settings                ChannelOverlayLayerSettings
	OverlayID               uuid.UUID
	PosX                    int
	PosY                    int
	Width                   int
	Height                  int
	Rotation                int
	CreatedAt               time.Time
	UpdatedAt               time.Time
	PeriodicallyRefetchData bool
}

type ChannelOverlay struct {
	ID        uuid.UUID
	ChannelID string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Width     int
	Height    int
	InstaSave bool
	Layers    []ChannelOverlayLayer
}

var ChannelOverlayNil = ChannelOverlay{}
