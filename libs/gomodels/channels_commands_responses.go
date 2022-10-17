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

type ChannelsCommandsResponses struct {
	ID        string      `gorm:"primaryKey;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	Text      null.String `gorm:"column:text;type:TEXT;"                          json:"text"`
	CommandID string      `gorm:"column:commandId;type:TEXT;"                     json:"commandId"`
	Order     int         `gorm:"column:order";type:INT" json:"order"`
}

func (c *ChannelsCommandsResponses) TableName() string {
	return "channels_commands_responses"
}
