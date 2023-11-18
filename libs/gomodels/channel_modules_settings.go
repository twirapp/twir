package model

import (
	"encoding/json"

	"github.com/guregu/null"
)

type ChannelModulesSettings struct {
	ID        string      `gorm:"column:id;type:uuid"        json:"id"`
	Type      string      `gorm:"column:type;"               json:"type"`
	Settings  []byte      `gorm:"column:settings;type:jsonb" json:"settings"`
	ChannelId string      `gorm:"column:channelId;type:text" json:"channelId"`
	UserId    null.String `gorm:"column:userId;type:text"    json:"userId"`
}

func (ChannelModulesSettings) TableName() string {
	return "channels_modules_settings"
}

type UserYoutubeSettings struct {
	MaxRequests  uint32 `json:"maxRequests"`
	MinWatchTime uint64 `json:"minWatchTime"`
	MinMessages  uint32 `json:"minMessages"`
	// in hours
	MinFollowTime uint32 `json:"minFollowTime"`
}

type SongYoutubeSettings struct {
	MaxLength          uint32   `json:"maxLength"`
	MinViews           uint64   `json:"minViews"`
	AcceptedCategories []string `json:"acceptedCategories"`
}

type BlackListYoutubeSettings struct {
	UsersIds     []string `json:"usersIds"`
	SongsIds     []string `json:"songsIds"`
	ChannelsIds  []string `json:"channelsIds"`
	ArtistsNames []string `json:"artistsNames"`
	Words        []string `json:"words"`
}

func emptize(slice []string) []string {
	if slice == nil {
		return []string{}
	} else {
		return slice
	}
}

func (s *BlackListYoutubeSettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		BlackListYoutubeSettings{
			UsersIds:     emptize(s.UsersIds),
			SongsIds:     emptize(s.SongsIds),
			ChannelsIds:  emptize(s.ChannelsIds),
			ArtistsNames: emptize(s.ArtistsNames),
		},
	)
}

type YoutubeSettings struct {
	AcceptOnlyWhenOnline    bool                     `json:"acceptOnlyWhenOnline"`
	ChannelPointsRewardName string                   `json:"channelPointsRewardName"`
	MaxRequests             uint16                   `json:"maxRequests"`
	User                    UserYoutubeSettings      `json:"user"                    validate:"required"`
	Song                    SongYoutubeSettings      `json:"song"                    validate:"required"`
	BlackList               BlackListYoutubeSettings `json:"blacklist"               validate:"required"`
}

type EightBallSettings struct {
	Answers []string `validate:"required" json:"answers"`
	Enabled bool     `json:"enabled"`
}

type RussianRouletteSetting struct {
	Enabled               bool `json:"enabled"`
	CanBeUsedByModerators bool `json:"canBeUsedByModerator"`
	TimeoutSeconds        int  `json:"timeoutTime"`
	DecisionSeconds       int  `json:"decisionTime"`
	TumberSize            int  `json:"tumberSize"`
	ChargedBullets        int  `json:"chargedBullets"`

	InitMessage    string `json:"initMessage"`
	SurviveMessage string `json:"surviveMessage"`
	DeathMessage   string `json:"deathMessage"`
}

type ChatAlertsSettings struct {
	Followers        ChatAlertsFollowersSettings `json:"followers"`
	Raids            ChatAlertsRaids             `json:"raids"`
	Donations        ChatAlertsDonations         `json:"donations"`
	Subscribers      ChatAlertsSubscribers       `json:"subscribers"`
	Cheers           ChatAlertsCheers            `json:"cheers"`
	Redemptions      ChatAlertsRedemptions       `json:"redemptions"`
	FirstUserMessage ChatAlertsFirstUserMessage  `json:"firstUserMessage"`
	StreamOnline     ChatAlertsStreamOnline      `json:"streamOnline"`
	StreamOffline    ChatAlertsStreamOffline     `json:"streamOffline"`
	ChatCleared      ChatAlertsChatCleared       `json:"chatCleared"`
	Ban              ChatAlertsBan               `json:"ban"`
}

type ChatAlertsFollowersSettings struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsCountedMessage struct {
	Count int    `json:"count"`
	Text  string `json:"text"`
}

type ChatAlertsMessage struct {
	Text string `json:"text"`
}

type ChatAlertsRaids struct {
	Enabled  bool                       `json:"enabled"`
	Messages []ChatAlertsCountedMessage `json:"messages"`
	Cooldown int                        `json:"cooldown"`
}

type ChatAlertsDonations struct {
	Enabled  bool                       `json:"enabled"`
	Messages []ChatAlertsCountedMessage `json:"messages"`
	Cooldown int                        `json:"cooldown"`
}

type ChatAlertsSubscribers struct {
	Enabled  bool                       `json:"enabled"`
	Messages []ChatAlertsCountedMessage `json:"messages"`
	Cooldown int                        `json:"cooldown"`
}

type ChatAlertsCheers struct {
	Enabled  bool                       `json:"enabled"`
	Messages []ChatAlertsCountedMessage `json:"messages"`
	Cooldown int                        `json:"cooldown"`
}

type ChatAlertsRedemptions struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsFirstUserMessage struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsStreamOnline struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsStreamOffline struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsChatCleared struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsBan struct {
	Enabled           bool                       `json:"enabled"`
	Messages          []ChatAlertsCountedMessage `json:"messages"`
	IgnoreTimeoutFrom []string                   `json:"ignoreTimeoutFrom"`
	Cooldown          int                        `json:"cooldown"`
}

type ChatOverlaySettings struct {
	MessageHideTimeout uint32 `json:"messageHideTimeout,omitempty"`
	MessageShowDelay   uint32 `json:"messageShowDelay,omitempty"`
	Preset             string `json:"preset,omitempty"`
	FontSize           uint32 `json:"fontSize,omitempty"`
	HideCommands       bool   `json:"hideCommands,omitempty"`
	HideBots           bool   `json:"hideBots,omitempty"`
	FontFamily         string `json:"fontFamily,omitempty"`
	ShowBadges         bool   `json:"showBadges,omitempty"`
	ShowAnnounceBadge  bool   `json:"showAnnounceBadge,omitempty"`
	ReverseMessages    bool   `json:"reverseMessages,omitempty"`
	TextShadowColor    string `json:"textShadowColor,omitempty"`
	TextShadowSize     uint32 `json:"textShadowSize,omitempty"`
}

type KappagenOverlaySettingsEmotes struct {
	Time  int32 `json:"time,omitempty"`
	Max   int32 `json:"max,omitempty"`
	Queue int32 `json:"queue,omitempty"`
}

type KappagenOverlaySettingsSize struct {
	// from 7 to 20
	RatioNormal int32 `json:"ratioNormal,omitempty"`
	// from 14 to 40
	RatioSmall int32 `json:"ratioSmall,omitempty"`
	Min        int32 `json:"min,omitempty"`
	Max        int32 `json:"max,omitempty"`
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
