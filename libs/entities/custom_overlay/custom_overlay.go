package customoverlayentity

import (
	"time"

	"github.com/google/uuid"
)

type ChannelOverlayType string

const (
	ChannelOverlayTypeHTML    ChannelOverlayType = "HTML"
	ChannelOverlayTypeIMAGE   ChannelOverlayType = "IMAGE"
	ChannelOverlayTypeTEXT    ChannelOverlayType = "TEXT"
	ChannelOverlayTypeVIDEO   ChannelOverlayType = "VIDEO"
	ChannelOverlayTypeIFRAME  ChannelOverlayType = "IFRAME"
	ChannelOverlayTypeYOUTUBE ChannelOverlayType = "YOUTUBE"
	ChannelOverlayTypeEMOTE   ChannelOverlayType = "EMOTE"
)

type ChannelOverlayLayerSettings struct {
	HtmlOverlayHTML                    string
	HtmlOverlayCSS                     string
	HtmlOverlayJS                      string
	HtmlOverlayDataPollSecondsInterval int
	ImageUrl                           string

	TextContent    string
	TextFontFamily string
	TextFontSize   int
	TextFontWeight int
	TextColor      string
	TextAlign      string

	VideoUrl   string
	VideoLoop  bool
	VideoMuted bool

	IframeUrl   string
	IframeScale float64
	WidgetKey   string

	YoutubeVideoID  string
	YoutubeAutoplay bool
	YoutubeLoop     bool
	YoutubeMuted    bool

	EmoteUrl      string
	EmoteName     string
	EmoteProvider string
}

type ChannelOverlayLayer struct {
	ID                      uuid.UUID
	Type                    ChannelOverlayType
	Name                    string
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
	Locked                  bool
	Visible                 bool
	Opacity                 float64
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

	isNil bool
}

func (c ChannelOverlay) IsNil() bool {
	return c.isNil
}

var ChannelOverlayNil = ChannelOverlay{
	isNil: true,
}
