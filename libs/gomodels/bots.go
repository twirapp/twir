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

type Bots struct {
	ID       string      `gorm:"primary_key;column:id;type:TEXT;" json:"id"`
	Type     string      `gorm:"column:type;type:VARCHAR;"        json:"type"`
	TokenID  null.String `gorm:"column:tokenId;type:TEXT;"        json:"tokenId"`
	Token    *Tokens     `gorm:"foreignKey:TokenID"               json:"token"`
	Channels []Channels  `gorm:"foreignKey:BotID"                 json:"channels"`
}

func (b *Bots) TableName() string {
	return "bots"
}
