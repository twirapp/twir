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

type ChannelsIntegrations struct {
	ID            string        `gorm:"primary_key;column:id;type:TEXT;"        json:"id"`
	Enabled       bool          `gorm:"column:enabled;type:BOOL;default:false;" json:"enabled"`
	ChannelID     string        `gorm:"column:channelId;type:TEXT;"             json:"channelId"`
	IntegrationID string        `gorm:"column:integrationId;type:TEXT;"         json:"integrationId"`
	AccessToken   null.String   `gorm:"column:accessToken;type:TEXT;"           json:"accessToken"`
	RefreshToken  null.String   `gorm:"column:refreshToken;type:TEXT;"          json:"refreshToken"`
	ClientID      null.String   `gorm:"column:clientId;type:TEXT;"              json:"clientId"`
	ClientSecret  null.String   `gorm:"column:clientSecret;type:TEXT;"          json:"clientSecret"`
	APIKey        null.String   `gorm:"column:apiKey;type:TEXT;"                json:"apiKey"`
	Data          null.String   `gorm:"column:data;type:JSONB;"                 json:"data"`
	Integration   *Integrations `gorm:"foreignKey:IntegrationID"                json:"-"`
}

func (c *ChannelsIntegrations) TableName() string {
	return "channels_integrations"
}
