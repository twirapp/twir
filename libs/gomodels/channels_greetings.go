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

type ChannelsGreetings struct {
	ID        string `gorm:"primaryKey;AUTO_INCREMENT;column:id;type:TEXT;default:uuid_generate_v4()" json:"id"`
	ChannelID string `gorm:"column:channelId;type:TEXT;"                     json:"channelId"`
	UserID    string `gorm:"column:userId;type:TEXT;"                        json:"userId"`
	Enabled   bool   `gorm:"column:enabled;type:bool;"                       json:"enabled"`
	Text      string `gorm:"column:text;type:TEXT;"                          json:"text"`
	IsReply   bool   `gorm:"column:isReply;type:bool"                        json:"isReply"`
	Processed bool   `gorm:"column:processed;type:bool"                      json:"processed"`
}

func (c *ChannelsGreetings) TableName() string {
	return "channels_greetings"
}
