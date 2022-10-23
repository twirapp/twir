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
	IsTester   bool           `gorm:"column:isTester;type:BOOL;default:false;"   json:"isTester"`
	IsBotAdmin bool           `gorm:"column:isBotAdmin;type:BOOL;default:false;" json:"isBotAdmin"`
	ApiKey     string         `gorm:"column:apiKey;type:TEXT;"                   json:"apiKey"`
}

func (u *Users) TableName() string {
	return "users"
}
