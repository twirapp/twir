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

type ChatAlertsSettings struct {
	Followers           ChatAlertsFollowersSettings   `json:"followers"`
	Raids               ChatAlertsRaids               `json:"raids"`
	Donations           ChatAlertsDonations           `json:"donations"`
	Subscribers         ChatAlertsSubscribers         `json:"subscribers"`
	Cheers              ChatAlertsCheers              `json:"cheers"`
	Redemptions         ChatAlertsRedemptions         `json:"redemptions"`
	FirstUserMessage    ChatAlertsFirstUserMessage    `json:"firstUserMessage"`
	StreamOnline        ChatAlertsStreamOnline        `json:"streamOnline"`
	StreamOffline       ChatAlertsStreamOffline       `json:"streamOffline"`
	ChatCleared         ChatAlertsChatCleared         `json:"chatCleared"`
	Ban                 ChatAlertsBan                 `json:"ban"`
	UnbanRequestCreate  ChatAlertsUnbanRequestCreate  `json:"unbanRequestCreate"`
	UnbanRequestResolve ChatAlertsUnbanRequestResolve `json:"unbanRequestResolve"`
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

type ChatAlertsUnbanRequestCreate struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}

type ChatAlertsUnbanRequestResolve struct {
	Enabled  bool                `json:"enabled"`
	Messages []ChatAlertsMessage `json:"messages"`
	Cooldown int                 `json:"cooldown"`
}
