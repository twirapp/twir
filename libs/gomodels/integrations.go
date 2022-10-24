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

type Integrations struct {
	ID           string         `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	Service      string         `gorm:"column:service;type:VARCHAR;"                    json:"service"`
	AccessToken  sql.NullString `gorm:"column:accessToken;type:TEXT;"                   json:"access_token"`
	RefreshToken sql.NullString `gorm:"column:refreshToken;type:TEXT;"                  json:"refresh_token"`
	ClientID     sql.NullString `gorm:"column:clientId;type:TEXT;"                      json:"client_id"`
	ClientSecret sql.NullString `gorm:"column:clientSecret;type:TEXT;"                  json:"client_secret"`
	APIKey       sql.NullString `gorm:"column:apiKey;type:TEXT;"                        json:"api_key"`
	RedirectURL  sql.NullString `gorm:"column:redirectUrl;type:TEXT;"                   json:"redirect_url"`
}

func (i *Integrations) TableName() string {
	return "integrations"
}
