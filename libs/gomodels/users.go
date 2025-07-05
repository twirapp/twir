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

type Users struct {
	ID         string         `gorm:"primary_key;column:id;type:TEXT;"           json:"id"`
	TokenID    sql.NullString `gorm:"column:tokenId;type:TEXT;"                  json:"tokenId"`
	IsBotAdmin bool           `gorm:"column:isBotAdmin;type:BOOL;default:false;" json:"isBotAdmin"`
	ApiKey     string         `gorm:"column:apiKey;type:TEXT;"                   json:"apiKey"`
	Channel    *Channels      `gorm:"foreignKey:ID"                              json:"channel"`
	Token      *Tokens        `gorm:"foreignKey:TokenID"                         json:"token"`
	Stats      *UsersStats    `gorm:"foreignKey:UserID"                          json:"stats"`
	IsBanned   bool           `gorm:"column:is_banned;type:BOOL;"       json:"isBanned"`
	CreatedAt  time.Time      `gorm:"column:created_at;type:TIMESTAMPTZ;default:now()" json:"createdAt"`

	HideOnLandingPage bool `gorm:"column:hide_on_landing_page;type:BOOL;default:false;" json:"hide_on_landing_page"`

	Roles []ChannelRoleUser `gorm:"foreignKey:UserID" json:"roles"`
}

func (u *Users) TableName() string {
	return "users"
}
