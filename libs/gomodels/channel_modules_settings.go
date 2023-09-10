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
