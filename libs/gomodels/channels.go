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

type Channels struct {
	ID             string `gorm:"primary_key;column:id;type:TEXT;"               json:"id"`
	IsEnabled      bool   `gorm:"column:isEnabled;type:BOOL;default:true;"       json:"isEnabled"`
	IsTwitchBanned bool   `gorm:"column:isTwitchBanned;type:BOOL;default:false;" json:"isTwitchBanned"`
	IsBanned       bool   `gorm:"column:isBanned;type:BOOL;default:false;"       json:"isBanned"`
	BotID          string `gorm:"column:botId;type:TEXT;"                        json:"botId"`
}

func (c *Channels) TableName() string {
	return "channels"
}
