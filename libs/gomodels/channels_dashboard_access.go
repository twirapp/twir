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

type ChannelsDashboardAccess struct {
	ID        string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	ChannelID string `gorm:"column:channelId;type:TEXT;"                     json:"channelId"`
	UserID    string `gorm:"column:userId;type:TEXT;"                        json:"userId"`
}

func (c *ChannelsDashboardAccess) TableName() string {
	return "channels_dashboard_access"
}
