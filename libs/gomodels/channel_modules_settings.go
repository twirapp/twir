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

//type JSON json.RawMessage
//func (j *JSON) Scan(value interface{}) error {
//	bytes, ok := value.([]byte)
//	if !ok {
//		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
//	}
//
//	result := json.RawMessage{}
//	err := json.Unmarshal(bytes, &result)
//	*j = JSON(result)
//	return err
//}
//
//// Value return json value, implement driver.Valuer interface
//func (j JSON) Value() (driver.Value, error) {
//	if len(j) == 0 {
//		return nil, nil
//	}
//
//	return json.RawMessage(j).MarshalJSON()
//}

type UserYoutubeSettings struct {
	MaxRequests  uint32 `validate:"required,lte=4294967295"           json:"maxRequests"`
	MinWatchTime uint64 `validate:"required,lte=18446744073709551615" json:"minWatchTime"`
	MinMessages  uint32 `validate:"required,lte=4294967295"           json:"minMessages"`
	// in hours
	MinFollowTime uint32 `validate:"required,lte=4294967295"           json:"minFollowTime"`
}

type SongYoutubeSettings struct {
	MaxLength          uint32   `validate:"required,lte=4294967295"           json:"maxLength"`
	MinViews           uint64   `validate:"required,lte=18446744073709551615" json:"minViews"`
	AcceptedCategories []string `                                             json:"acceptedCategories"`
}

type BlackListYoutubeSettings struct {
	UsersIds     []string `validate:"required" json:"usersIds"`
	SongsIds     []string `validate:"required" json:"songsIds"`
	ChannelsIds  []string `validate:"required" json:"channelsIds"`
	ArtistsNames []string `validate:"required" json:"artistsNames"`
}

func emptize(slice []string) []string {
	if slice == nil {
		return []string{}
	} else {
		return slice
	}
}

func (s *BlackListYoutubeSettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(BlackListYoutubeSettings{
		UsersIds:     emptize(s.UsersIds),
		SongsIds:     emptize(s.SongsIds),
		ChannelsIds:  emptize(s.ChannelsIds),
		ArtistsNames: emptize(s.ArtistsNames),
	})
}

type YoutubeSettings struct {
	AcceptOnlyWhenOnline    bool                     `json:"acceptOnlyWhenOnline"`
	ChannelPointsRewardName string                   `json:"channelPointsRewardName"`
	MaxRequests             uint16                   `json:"maxRequests"`
	User                    UserYoutubeSettings      `json:"user"                    validate:"required"`
	Song                    SongYoutubeSettings      `json:"song"                    validate:"required"`
	BlackList               BlackListYoutubeSettings `json:"blacklist"               validate:"required"`
}
