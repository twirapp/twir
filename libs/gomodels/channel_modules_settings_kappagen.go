package model

type KappagenOverlaySettingsEmotesEmojiStyle int32

const (
	KappagenOverlaySettingsEmotesEmojiStyle_None     KappagenOverlaySettingsEmotesEmojiStyle = 0
	KappagenOverlaySettingsEmotesEmojiStyle_Twemoji  KappagenOverlaySettingsEmotesEmojiStyle = 1
	KappagenOverlaySettingsEmotesEmojiStyle_Openmoji KappagenOverlaySettingsEmotesEmojiStyle = 2
	KappagenOverlaySettingsEmotesEmojiStyle_Noto     KappagenOverlaySettingsEmotesEmojiStyle = 3
	KappagenOverlaySettingsEmotesEmojiStyle_Blobmoji KappagenOverlaySettingsEmotesEmojiStyle = 4
)

type KappagenOverlaySettingsEmotes struct {
	Time           int32                                   `json:"time,omitempty"`
	Max            int32                                   `json:"max,omitempty"`
	Queue          int32                                   `json:"queue,omitempty"`
	FfzEnabled     bool                                    `json:"ffzEnabled,omitempty"`
	BttvEnabled    bool                                    `json:"bttvEnabled,omitempty"`
	SevenTvEnabled bool                                    `json:"sevenTvEnabled,omitempty"`
	EmojiStyle     KappagenOverlaySettingsEmotesEmojiStyle `json:"emojiStyle,omitempty"`
}

type KappagenOverlaySettingsSize struct {
	// from 7 to 20
	RatioNormal float64 `json:"ratioNormal,omitempty"`
	// from 14 to 40
	RatioSmall float64 `json:"ratioSmall,omitempty"`
	Min        int32   `json:"min,omitempty"`
	Max        int32   `json:"max,omitempty"`
}

type KappagenOverlaySettingsCube struct {
	Speed int32 `json:"speed,omitempty"`
}

type KappagenOverlaySettingsAnimation struct {
	FadeIn  bool `json:"fadeIn,omitempty"`
	FadeOut bool `json:"fadeOut,omitempty"`
	ZoomIn  bool `json:"zoomIn,omitempty"`
	ZoomOut bool `json:"zoomOut,omitempty"`
}

type KappagenOverlaySettingsAnimationSettingsPrefs struct {
	Size    *float64 `json:"size"`
	Center  *bool    `json:"center"`
	Speed   *int32   `json:"speed"`
	Faces   *bool    `json:"faces"`
	Message []string `json:"message"`
	Time    *int32   `json:"time"`
}

type KappagenOverlaySettingsAnimationSettings struct {
	Style   string                                         `json:"style"`
	Prefs   *KappagenOverlaySettingsAnimationSettingsPrefs `json:"prefs"`
	Count   *int32                                         `json:"count"`
	Enabled bool                                           `json:"enabled"`
}

type KappagenOverlaySettingsEvent struct {
	Event          int32    `json:"event"`
	DisabledStyles []string `json:"disabledStyles,omitempty"`
	Enabled        bool     `json:"enabled,omitempty"`
}

type KappagenOverlaySettings struct {
	Emotes      KappagenOverlaySettingsEmotes              `json:"emotes,omitempty"`
	Size        KappagenOverlaySettingsSize                `json:"size,omitempty"`
	Cube        KappagenOverlaySettingsCube                `json:"cube,omitempty"`
	Animation   KappagenOverlaySettingsAnimation           `json:"animation,omitempty"`
	Animations  []KappagenOverlaySettingsAnimationSettings `json:"animations,omitempty"`
	EnableRave  bool                                       `json:"enableRave,omitempty"`
	Events      []KappagenOverlaySettingsEvent             `json:"events,omitempty"`
	EnableSpawn bool                                       `json:"enableSpawn,omitempty"`
}
