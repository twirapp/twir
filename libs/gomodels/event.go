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

type Event struct {
	ID        string      `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	ChannelID string      `gorm:"column:channelId;type:TEXT;" json:"channelId"`
	Type      string      `gorm:"column:type;type:TEXT;"                     json:"type"`
	RewardID  null.String `gorm:"column:rewardId;type:TEXT;"                     json:"rewardId"`
	CommandID null.String `gorm:"column:commandId;type:TEXT;"                     json:"commandId"`

	Operations []EventOperation
}

func (c *Event) TableName() string {
	return "channels_events"
}
