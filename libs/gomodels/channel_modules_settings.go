package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JSON json.RawMessage

type ChannelModulesSettings struct {
	ID       string `gorm:"column:id;type:uuid"     json:"id"`
	Type     string `gorm:"column:type;type:string" json:"type"`
	Settings JSON   `gorm:"column:settings"         json:"settings"`
}

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

type UserYoutubeSettings struct {
	MaxRequests   *uint `json:"maxRequests"`
	MinWatchTime  *uint `json:"minWatchTime"`
	MinMessages   *uint `json:"minMessages"`
	MinFollowTime *uint `json:"minFollowTime"`
}

type SongYoutubeSettings struct {
	MaxLength *uint `json:"maxLength"`
	MinViews  *uint `json:"minViews"`
}

type BlackListYoutubeSettings struct {
	UsersIds     *[]string `json:"usersIds"`
	SongsIds     *[]string `json:"songsIds"`
	ChannelsIds  *[]string `json:"channelsIds"`
	ArtistsNames *[]string `json:"artistsNames"`
}

type YoutubeSettings struct {
	AcceptOnlyWhenOnline    *bool                     `json:"acceptOnlyWhenOnline"`
	ChannelPointsRewardName *string                   `json:"channelPointsRewardName"`
	User                    *UserYoutubeSettings      `json:"user"`
	Song                    *SongYoutubeSettings      `json:"song"`
	BlackList               *BlackListYoutubeSettings `json:"blaclist"`
}
