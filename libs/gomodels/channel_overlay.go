package model

import (
	"time"

	"github.com/google/uuid"
)

type ChannelOverlay struct {
	ID        uuid.UUID `gorm:"primary_key;column:id;type:UUID;"  json:"id"`
	ChannelID string    `gorm:"column:channel_id;type:TEXT;"  json:"channelId"`
	Name      string    `gorm:"column:name;type:TEXT;"  json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;data:timestamp;"  json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;data:timestamp;"  json:"updatedAt"`

	Channel *Channels `gorm:"foreignKey:ChannelID" json:"channel"`

	Layers []ChannelOverlayLayer `gorm:"foreignKey:OverlayID" json:"layers"`
}

func (c ChannelOverlay) TableName() string {
	return "channels_overlays"
}
