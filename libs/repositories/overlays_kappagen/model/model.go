package model

import (
	"time"

	"github.com/google/uuid"
	eventmodel "github.com/twirapp/twir/libs/repositories/events/model"
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
	ID        uuid.UUID               `json:"id,omitempty"`
	ChannelID string                  `json:"channel_id,omitempty"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
	Settings  KappagenOverlaySettings `json:"settings"`
}

type KappagenOverlaySettings struct {
	EnableSpawn    bool                                `json:"enable_spawn,omitempty"`
	ExcludedEmotes []string                            `json:"excluded_emotes,omitempty"`
	EnableRave     bool                                `json:"enable_rave,omitempty"`
	Animation      KappagenOverlayAnimationSettings    `json:"animation"`
	Animations     []KappagenOverlayAnimationsSettings `json:"animations,omitempty"`
	Emotes         KappagenOverlayEmotesSettings       `json:"emotes"`
	Size           KappagenOverlaySizeSettings         `json:"size"`
	Events         []KappagenOverlayEvent              `json:"events,omitempty,omitzero"`
}

type KappagenOverlayEmotesSettings struct {
	Time           int                `json:"time,omitempty"`
	Max            int                `json:"max,omitempty"`
	Queue          int                `json:"queue,omitempty"`
	FfzEnabled     bool               `json:"ffz_enabled,omitempty"`
	BttvEnabled    bool               `json:"bttv_enabled,omitempty"`
	SevenTvEnabled bool               `json:"seven_tv_enabled,omitempty"`
	EmojiStyle     KappagenEmojiStyle `json:"emoji_style,omitempty"`
}

type KappagenOverlaySizeSettings struct {
	RatioNormal float64 `json:"ratio_normal,omitempty"`
	RatioSmall  float64 `json:"ratio_small,omitempty"`
	Min         int     `json:"min,omitempty"`
	Max         int     `json:"max,omitempty"`
}

type KappagenOverlayAnimationSettings struct {
	FadeIn  bool `json:"fade_in,omitempty"`
	FadeOut bool `json:"fade_out,omitempty"`
	ZoomIn  bool `json:"zoom_in,omitempty"`
	ZoomOut bool `json:"zoom_out,omitempty"`
}

type KappagenOverlayAnimationsPrefsSettings struct {
	Size    float64  `json:"size,omitempty"`
	Center  bool     `json:"center,omitempty"`
	Speed   int      `json:"speed,omitempty"`
	Faces   bool     `json:"faces,omitempty"`
	Message []string `json:"message,omitempty"`
	Time    int      `json:"time,omitempty"`
}

type KappagenOverlayAnimationsSettings struct {
	Style   string                                 `json:"style,omitempty"`
	Prefs   KappagenOverlayAnimationsPrefsSettings `json:"prefs"`
	Count   int                                    `json:"count,omitempty"`
	Enabled bool                                   `json:"enabled,omitempty"`
}

type KappagenOverlayEvent struct {
	Event              eventmodel.EventType `json:"event,omitempty"`
	DisabledAnimations []string             `json:"disabled_animations,omitempty"`
	Enabled            bool                 `json:"enabled,omitempty"`
}

var Nil = KappagenOverlay{}
