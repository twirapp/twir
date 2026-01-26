package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type ChannelsTimers struct {
	ID                       string                     `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;"      json:"id"`
	ChannelID                string                     `gorm:"column:channelId;type:TEXT;"                          json:"-"`
	Name                     string                     `gorm:"column:name;type:VARCHAR;size:255;"                   json:"name"`
	Enabled                  bool                       `gorm:"column:enabled;type:BOOL;"                            json:"enabled"`
	OfflineEnabled           bool                       `gorm:"column:offline_enabled;type:BOOL;default:false;"       json:"offlineEnabled"`
	OnlineEnabled            bool                       `gorm:"column:online_enabled;type:BOOL;default:true;"         json:"onlineEnabled"`
	TimeInterval             int32                      `gorm:"column:timeInterval;type:INT4;default:0;"             json:"timeInterval"`
	MessageInterval          int32                      `gorm:"column:messageInterval;type:INT4;default:0;"          json:"messageInterval"`
	LastTriggerMessageNumber int32                      `gorm:"column:lastTriggerMessageNumber;type:INT4;default:0;" json:"-"`
	Responses                []*ChannelsTimersResponses `gorm:"foreignKey:TimerID"                                   json:"responses"`
	Channel                  *Channels                  `gorm:"foreignKey:ChannelID" json:"channel"`
	AnnounceColor            int                        `gorm:"column:announce_color;type:INT4;default:0;"           json:"announceColor"`
}

func (c *ChannelsTimers) TableName() string {
	return "channels_timers"
}
