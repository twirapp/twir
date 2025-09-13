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

type ChannelsKeywords struct {
	ID               string    `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	ChannelID        string    `gorm:"column:channelId;type:TEXT;"                     json:"channelId"`
	Text             string    `gorm:"column:text;type:TEXT;"                          json:"text"`
	Response         string    `gorm:"column:response;type:TEXT;"                      json:"response"`
	Enabled          bool      `gorm:"column:enabled;type:BOOL"                        json:"enabled"`
	Cooldown         int       `gorm:"column:cooldown;type:INT4;default:0;"            json:"cooldown"         swaggertype:"integer"`
	CooldownExpireAt null.Time `gorm:"column:cooldownExpireAt;type:timestamp;"         json:"cooldownExpireAt" swaggertype:"string"`
	IsReply          bool      `gorm:"column:isReply;type:BOOL"                        json:"isReply"`
	IsRegular        bool      `gorm:"column:isRegular;type:bool"                      json:"isRegular"`
	Usages           int       `gorm:"column:usages;type:int4"                         json:"usages"`
	RolesIDs         []string  `gorm:"roles_ids;type:text[];default:[];" json:"roles_ids"`
}

func (c *ChannelsKeywords) TableName() string {
	return "channels_keywords"
}
