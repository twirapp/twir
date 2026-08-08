package model

import (
	"time"

	"github.com/google/uuid"
)

type OverlayType string

const (
	OverlayTypeHTML    OverlayType = "HTML"
	OverlayTypeIMAGE   OverlayType = "IMAGE"
	OverlayTypeTEXT    OverlayType = "TEXT"
	OverlayTypeVIDEO   OverlayType = "VIDEO"
	OverlayTypeIFRAME  OverlayType = "IFRAME"
	OverlayTypeYOUTUBE OverlayType = "YOUTUBE"
	OverlayTypeEMOTE   OverlayType = "EMOTE"
)

type OverlayLayerSettings struct {
	HtmlOverlayHTML                    string `json:"htmlOverlayHtml"`
	HtmlOverlayCSS                     string `json:"htmlOverlayCss"`
	HtmlOverlayJS                      string `json:"htmlOverlayJs"`
	HtmlOverlayDataPollSecondsInterval int    `json:"htmlOverlayDataPollSecondsInterval"`
	ImageUrl                           string `json:"imageUrl"`

	TextContent    string `json:"textContent"`
	TextFontFamily string `json:"textFontFamily"`
	TextFontSize   int    `json:"textFontSize"`
	TextFontWeight int    `json:"textFontWeight"`
	TextColor      string `json:"textColor"`
	TextAlign      string `json:"textAlign"`

	VideoUrl   string `json:"videoUrl"`
	VideoLoop  bool   `json:"videoLoop"`
	VideoMuted bool   `json:"videoMuted"`

	IframeUrl   string  `json:"iframeUrl"`
	IframeScale float64 `json:"iframeScale"`
	WidgetKey   string  `json:"widgetKey"`

	YoutubeVideoID  string `json:"youtubeVideoId"`
	YoutubeAutoplay bool   `json:"youtubeAutoplay"`
	YoutubeLoop     bool   `json:"youtubeLoop"`
	YoutubeMuted    bool   `json:"youtubeMuted"`

	EmoteUrl      string `json:"emoteUrl"`
	EmoteName     string `json:"emoteName"`
	EmoteProvider string `json:"emoteProvider"`
}

type OverlayLayer struct {
	ID                      uuid.UUID            `json:"id"`
	Type                    OverlayType          `json:"type"`
	Name                    string               `json:"name"`
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
	Locked                  bool                 `json:"locked"`
	Visible                 bool                 `json:"visible"`
	Opacity                 float64              `json:"opacity"`

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
	InstaSave bool           `json:"insta_save"`
	Layers    []OverlayLayer `json:"layers"`

	isNil bool
}

func (o Overlay) IsNil() bool {
	return o.isNil
}

var Nil = Overlay{isNil: true}
