package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChannelOverlayType string

const (
	ChannelOverlayTypeHTML ChannelOverlayType = "HTML"
)

type ChannelOverlayLayerSettings struct {
	HtmlOverlayHTML                    string
	HtmlOverlayCSS                     string
	HtmlOverlayJS                      string
	HtmlOverlayDataPollSecondsInterval int
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
	Layers    []ChannelOverlayLayer
}

var ChannelOverlayNil = ChannelOverlay{}
