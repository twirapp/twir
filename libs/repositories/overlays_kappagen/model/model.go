package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
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
	CreatedAt      time.Time
	UpdatedAt      time.Time
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
	ID          ulid.ULID
	AnimationID ulid.ULID
	Size        float64
	Center      bool
	Speed       int
	Faces       bool
	Message     []string
	Time        int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type KappagenOverlayAnimationsSettings struct {
	ID        ulid.ULID
	OverlayID uuid.UUID
	Style     string
	Prefs     KappagenOverlayAnimationsPrefsSettings
	Count     int
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

var Nil = KappagenOverlay{}
