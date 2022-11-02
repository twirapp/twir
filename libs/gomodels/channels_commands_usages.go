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

type ChannelsCommandsUsages struct {
	ID        string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	UserID    string `gorm:"column:userId;type:TEXT;"                        json:"user_id"`
	ChannelID string `gorm:"column:channelId;type:TEXT;"                     json:"channel_id"`
	CommandID string `gorm:"column:commandId;type:TEXT;"                     json:"command_id"`
}

func (c *ChannelsCommandsUsages) TableName() string {
	return "channels_commands_usages"
}
