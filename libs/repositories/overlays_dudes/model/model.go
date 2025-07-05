package model

import (
	"time"

	"github.com/google/uuid"
)

type OverlaysDudes struct {
	ID        uuid.UUID
	ChannelID string
	CreatedAt time.Time
	
	// Dude settings
	DudeColor          string
	DudeEyesColor      string
	DudeCosmeticsColor string
	DudeMaxLifeTime    int
	DudeGravity        int
	DudeScale          float32
	DudeSoundsEnabled  bool
	DudeSoundsVolume   float32
	DudeVisibleName    bool
	DudeGrowTime       int
	DudeGrowMaxScale   int
	DudeMaxOnScreen    int
	DudeDefaultSprite  string
	
	// Message box settings
	MessageBoxEnabled      bool
	MessageBoxBorderRadius int
	MessageBoxBoxColor     string
	MessageBoxFontFamily   string
	MessageBoxFontSize     int
	MessageBoxPadding      int
	MessageBoxShowTime     int
	MessageBoxFill         string
	
	// Name box settings
	NameBoxFontFamily         string
	NameBoxFontSize           int
	NameBoxFill               []string
	NameBoxLineJoin           string
	NameBoxStrokeThickness    int
	NameBoxStroke             string
	NameBoxFillGradientStops  []float32
	NameBoxFillGradientType   int
	NameBoxFontStyle          string
	NameBoxFontVariant        string
	NameBoxFontWeight         int
	NameBoxDropShadow         bool
	NameBoxDropShadowAlpha    float32
	NameBoxDropShadowAngle    float32
	NameBoxDropShadowBlur     float32
	NameBoxDropShadowDistance float32
	NameBoxDropShadowColor    string
	
	// Ignore settings
	IgnoreCommands bool
	IgnoreUsers    bool
	IgnoredUsers   []string
	
	// Spitter emote settings
	SpitterEmoteEnabled bool
}

var Nil = OverlaysDudes{}
