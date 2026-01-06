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
	ID             string    `gorm:"primaryKey;column:id;type:TEXT;"               json:"id"`
	IsEnabled      bool      `gorm:"column:isEnabled;type:BOOL;"       json:"isEnabled"`
	IsTwitchBanned bool      `gorm:"column:isTwitchBanned;type:BOOL;" json:"isTwitchBanned"`
	IsBotMod       bool      `gorm:"column:isBotMod;type:BOOL;" json:"isBotMod"`
	BotID          string    `gorm:"column:botId;type:TEXT;"                        json:"botId"`
	PlanID         *string   `gorm:"column:plan_id;type:UUID;"                      json:"planId"`
	CreatedAt      time.Time `gorm:"column:created_at;type:TIMESTAMPTZ;default:now()" json:"createdAt"`

	Commands []ChannelsCommands `gorm:"foreignKey:ChannelID" json:"commands"`
	Roles    []*ChannelRole     `gorm:"foreignKey:ChannelID" json:"roles"`
	User     *Users             `gorm:"foreignKey:ID" json:"user"`
}

func (c *Channels) TableName() string {
	return "channels"
}
