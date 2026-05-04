package model

import (
	"time"

	"github.com/twirapp/twir/libs/entities/platform"
)

type ChannelInfoHistory struct {
	ID        string            `gorm:"primary_key;column:id;type:TEXT;" json:"id"`
	Category  string            `gorm:"column:category;type:text"               json:"category"`
	Title     string            `gorm:"column:title;type:text"               json:"title"`
	CreatedAt time.Time         `gorm:"column:createdAt;type:timestamp"         json:"createdAt"`
	ChannelID string            `gorm:"column:channelId;type:TEXT;" json:"channelId"`
	Platform  platform.Platform `gorm:"column:platform;type:platform" json:"platform"`

	Channel *Channels `gorm:"foreignKey:ID" json:"channel"`
}

func (c *ChannelInfoHistory) TableName() string {
	return "channels_info_history"
}
