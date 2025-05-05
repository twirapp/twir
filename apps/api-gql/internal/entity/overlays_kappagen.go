package entity

import (
	"github.com/google/uuid"
)

type KappagenEmojiStyle string

const (
	KappagenEmojiStyleOpenmoji KappagenEmojiStyle = "OPENMOJI"
	KappagenEmojiStyleTwemoji  KappagenEmojiStyle = "TWEMOJI"
	KappagenEmojiStyleBlobmoji KappagenEmojiStyle = "BLOBMOJI"
	KappagenEmojiStyleNoto     KappagenEmojiStyle = "NOTO"
)

type KappagenOverlay struct {
	ID             uuid.UUID                           `json:"id"`
	EnableSpawn    bool                                `json:"enableSpawn"`
	ExcludedEmotes []string                            `json:"excludedEmotes"`
	EnableRave     bool                                `json:"enableRave"`
	Animation      KappagenOverlayAnimationSettings    `json:"animation"`
	Animations     []KappagenOverlayAnimationsSettings `json:"animations"`
}

type KappagenOverlayEmotesSettings struct {
	Time           int                `json:"time"`
	Max            int                `json:"max"`
	Queue          int                `json:"queue"`
	FfzEnabled     bool               `json:"ffz_enabled"`
	BttvEnabled    bool               `json:"bttv_enabled"`
	SevenTvEnabled bool               `json:"seven_tv_enabled"`
	EmojiStyle     KappagenEmojiStyle `json:"emoji_style"`
}

type KappagenOverlaySizeSettings struct {
	RationNormal float64 `json:"rationNormal"`
	RationSmall  float64 `json:"rationSmall"`
	Min          int     `json:"min"`
	Max          int     `json:"max"`
}

type KappagenOverlayCubeSettings struct {
	Speed int `json:"speed"`
}

type KappagenOverlayAnimationSettings struct {
	FadeIn  bool `json:"fadeIn"`
	FadeOut bool `json:"fadeOut"`
	ZoomIn  bool `json:"zoomIn"`
	ZoomOut bool `json:"zoomOut"`
}

type KappagenOverlayAnimationsPrefsSettings struct {
	Size    float64  `json:"size"`
	Center  bool     `json:"center"`
	Speed   int      `json:"speed"`
	Faces   bool     `json:"faces"`
	Message []string `json:"message"`
	Time    int      `json:"time"`
}

type KappagenOverlayAnimationsSettings struct {
	Style   string                                 `json:"style"`
	Prefs   KappagenOverlayAnimationsPrefsSettings `json:"prefs"`
	Count   int                                    `json:"count"`
	Enabled bool                                   `json:"enabled"`
}
