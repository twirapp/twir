package entity

import (
	"time"

	"github.com/google/uuid"
)

type KappagenEmojiStyle int

const (
	KappagenEmojiStyleNone KappagenEmojiStyle = iota
	KappagenEmojiStyleTwemoji
	KappagenEmojiStyleOpenmoji
	KappagenEmojiStyleNoto
	KappagenEmojiStyleBlobmoji
)

type KappagenOverlay struct {
	ID        uuid.UUID
	ChannelID string
	CreatedAt time.Time
	UpdatedAt time.Time
	Settings  KappagenOverlaySettings
}

type KappagenOverlaySettings struct {
	EnableSpawn    bool
	ExcludedEmotes []string
	EnableRave     bool
	Animation      KappagenOverlayAnimationSettings
	Animations     []KappagenOverlayAnimationsSettings
	Emotes         KappagenOverlayEmotesSettings
	Size           KappagenOverlaySizeSettings
	Events         []KappagenOverlayEvent
}

type KappagenOverlayEmotesSettings struct {
	Time           int
	Max            int
	Queue          int
	FfzEnabled     bool
	BttvEnabled    bool
	SevenTvEnabled bool
	EmojiStyle     KappagenEmojiStyle
}

type KappagenOverlaySizeSettings struct {
	RatioNormal float64
	RatioSmall  float64
	Min         int
	Max         int
}

type KappagenOverlayAnimationSettings struct {
	FadeIn  bool
	FadeOut bool
	ZoomIn  bool
	ZoomOut bool
}

type KappagenOverlayAnimationsPrefsSettings struct {
	Size    float64
	Center  bool
	Speed   int
	Faces   bool
	Message []string
	Time    int
}

type KappagenOverlayAnimationsSettings struct {
	Style   string
	Prefs   KappagenOverlayAnimationsPrefsSettings
	Count   int
	Enabled bool
}

type KappagenOverlayEvent struct {
	Event              EventType
	DisabledAnimations []string
	Enabled            bool
}
