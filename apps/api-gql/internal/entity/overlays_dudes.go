package entity

import (
	"time"

	"github.com/google/uuid"
)

type DudesDudeSettings struct {
	Color          string
	EyesColor      string
	CosmeticsColor string
	MaxLifeTime    int
	Gravity        int
	Scale          float64
	SoundsEnabled  bool
	SoundsVolume   float64
	VisibleName    bool
	GrowTime       int
	GrowMaxScale   int
	MaxOnScreen    int
	DefaultSprite  string
}

type DudesMessageBoxSettings struct {
	Enabled      bool
	BorderRadius int
	BoxColor     string
	FontFamily   string
	FontSize     int
	Padding      int
	ShowTime     int
	Fill         string
}

type DudesNameBoxSettings struct {
	FontFamily         string
	FontSize           int
	Fill               []string
	LineJoin           string
	StrokeThickness    int
	Stroke             string
	FillGradientStops  []float64
	FillGradientType   int
	FontStyle          string
	FontVariant        string
	FontWeight         int
	DropShadow         bool
	DropShadowAlpha    float64
	DropShadowAngle    float64
	DropShadowBlur     float64
	DropShadowDistance float64
	DropShadowColor    string
}

type DudesIgnoreSettings struct {
	IgnoreCommands bool
	IgnoreUsers    bool
	Users          []string
}

type DudesSpitterEmoteSettings struct {
	Enabled bool
}

type DudesOverlaySettings struct {
	ID                   uuid.UUID
	DudeSettings         DudesDudeSettings
	MessageBoxSettings   DudesMessageBoxSettings
	NameBoxSettings      DudesNameBoxSettings
	IgnoreSettings       DudesIgnoreSettings
	SpitterEmoteSettings DudesSpitterEmoteSettings
	CreatedAt            time.Time
}

var DudesOverlaySettingsNil = DudesOverlaySettings{}
