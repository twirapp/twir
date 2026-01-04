package model

import (
	"time"

	"github.com/google/uuid"
)

type OverlayType string

const (
	OverlayTypeHTML  OverlayType = "HTML"
	OverlayTypeIMAGE OverlayType = "IMAGE"
)

type OverlayLayerSettings struct {
	HtmlOverlayHTML                    string `json:"htmlOverlayHtml"`
	HtmlOverlayCSS                     string `json:"htmlOverlayCss"`
	HtmlOverlayJS                      string `json:"htmlOverlayJs"`
	HtmlOverlayDataPollSecondsInterval int    `json:"htmlOverlayDataPollSecondsInterval"`
	ImageUrl                           string `json:"imageUrl"`
}

type OverlayLayer struct {
	ID                      uuid.UUID            `json:"id"`
	Type                    OverlayType          `json:"type"`
	Settings                OverlayLayerSettings `json:"settings"`
	OverlayID               uuid.UUID            `json:"overlay_id"`
	PosX                    int                  `json:"pos_x"`
	PosY                    int                  `json:"pos_y"`
	Width                   int                  `json:"width"`
	Height                  int                  `json:"height"`
	Rotation                int                  `json:"rotation"`
	CreatedAt               time.Time            `json:"created_at"`
	UpdatedAt               time.Time            `json:"updated_at"`
	PeriodicallyRefetchData bool                 `json:"periodically_refetch_data"`

	isNil bool
}

func (o OverlayLayer) IsNil() bool {
	return o.isNil
}

var LayerNil = OverlayLayer{isNil: true}

type Overlay struct {
	ID        uuid.UUID      `json:"id"`
	ChannelID string         `json:"channel_id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Width     int            `json:"width"`
	Height    int            `json:"height"`
	Layers    []OverlayLayer `json:"layers"`

	isNil bool
}

func (o Overlay) IsNil() bool {
	return o.isNil
}

var Nil = Overlay{isNil: true}
