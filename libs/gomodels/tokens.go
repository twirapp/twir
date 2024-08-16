package model

import (
	"database/sql"
	"time"

	"github.com/lib/pq"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type Tokens struct {
	ID                  string         `gorm:"column:id;type:TEXT;" json:"id"`
	AccessToken         string         `gorm:"column:accessToken;type:TEXT;"                   json:"accessToken"`
	RefreshToken        string         `gorm:"column:refreshToken;type:TEXT;"                  json:"refreshToken"`
	ExpiresIn           int32          `gorm:"column:expiresIn;type:INT4;"                     json:"expiresIn"`
	ObtainmentTimestamp time.Time      `gorm:"column:obtainmentTimestamp;type:TIMESTAMP;"      json:"obtainmentTimestamp"`
	Scopes              pq.StringArray `gorm:"column:scopes;type:text[]" json:"scopes"`
}

// TableName sets the insert table name for this struct type
func (t *Tokens) TableName() string {
	return "tokens"
}
