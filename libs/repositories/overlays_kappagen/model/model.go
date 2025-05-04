package model

import (
	"github.com/google/uuid"
)

type KappagenEmojiStyle int32

const (
	KappagenEmojiStyleNone     KappagenEmojiStyle = 0
	KappagenEmojiStyleTwemoji  KappagenEmojiStyle = 1
	KappagenEmojiStyleOpenmoji KappagenEmojiStyle = 2
	KappagenEmojiStyleNoto     KappagenEmojiStyle = 3
	KappagenEmojiStyleBlobmoji KappagenEmojiStyle = 4
)

type KappagenOverlay struct {
	ID             uuid.UUID
	ChannelID      string
	EnableSpawn    bool
	ExcludedEmotes []string
	EnableRave     bool
	Animation      KappagenOverlayAnimationSettings
	Animations     []KappagenOverlayAnimationsSettings
	Emotes         KappagenOverlayEmotesSettings
	Size           KappagenOverlaySizeSettings
	Cube           KappagenOverlayCubeSettings
}

type KappagenOverlayEmotesSettings struct {
	Time          int32
	Max           int32
	Queue         int32
	FfzEnabled    bool
	BttvEnabled   bool
	SevenTvEnabled bool
	EmojiStyle    KappagenEmojiStyle
}

type KappagenOverlaySizeSettings struct {
	RatioNormal float64
	RatioSmall  float64
	Min         int32
	Max         int32
}

type KappagenOverlayCubeSettings struct {
	Speed int32
}

type KappagenOverlayAnimationSettings struct {
	FadeIn  bool
	FadeOut bool
	ZoomIn  bool
	ZoomOut bool
}

type KappagenOverlayAnimationsPrefsSettings struct {
	Size    *float64
	Center  *bool
	Speed   *int32
	Faces   *bool
	Message []string
	Time    *int32
}

type KappagenOverlayAnimationsSettings struct {
	Style   string
	Prefs   *KappagenOverlayAnimationsPrefsSettings
	Count   *int32
	Enabled bool
}

var Nil = KappagenOverlay{}
