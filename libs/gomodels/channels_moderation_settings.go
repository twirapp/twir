package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type ChannelsModerationSettings struct {
	ID                 string         `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;"  json:"id"`
	Type               string         `gorm:"column:type;type:VARCHAR;"                        json:"type"`
	ChannelID          string         `gorm:"column:channelId;type:TEXT;"                      json:"channelId"`
	Enabled            bool           `gorm:"column:enabled;type:BOOL;default:false;"          json:"enabled"`
	Subscribers        bool           `gorm:"column:subscribers;type:BOOL;default:false;"      json:"subscribers"`
	Vips               bool           `gorm:"column:vips;type:BOOL;default:false;"             json:"vips"`
	BanTime            int32          `gorm:"column:banTime;type:INT4;default:600;"            json:"banTime"`
	BanMessage         string         `gorm:"column:banMessage;type:TEXT;"                     json:"banMessage"         swaggertype:"string"`
	WarningMessage     string         `gorm:"column:warningMessage;type:TEXT;"                 json:"warningMessage"     swaggertype:"string"`
	CheckClips         null.Bool      `gorm:"column:checkClips;type:BOOL;default:false;"       json:"checkClips"         swaggertype:"boolean"`
	TriggerLength      null.Int       `gorm:"column:triggerLength;type:INT4;default:300;"      json:"triggerLength"      swaggertype:"integer"`
	MaxPercentage      null.Int       `gorm:"column:maxPercentage;type:INT4;default:50;"       json:"maxPercentage"      swaggertype:"integer"`
	BlackListSentences pq.StringArray `gorm:"column:blackListSentences;type:JSONB;default:[];" json:"blackListSentences"`
}

func (c *ChannelsModerationSettings) TableName() string {
	return "channels_moderation_settings"
}
