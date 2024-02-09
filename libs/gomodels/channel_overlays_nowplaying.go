package model

import (
	"github.com/google/uuid"
	"github.com/satont/twir/libs/types/types/api/overlays"
)

type ChannelOverlayNowPlaying struct {
	ID     uuid.UUID                               `gorm:"primary_key;column:id;type:UUID;"  json:"id"`
	Preset overlays.ChannelOverlayNowPlayingPreset `gorm:"primary_key;column:preset;type:varchar;"  json:"preset"`

	ChannelID string    `gorm:"column:channel_id;type:text" json:"channel_id"`
	Channel   *Channels `gorm:"foreignkey:ChannelID;" json:"channel"`
}

func (c ChannelOverlayNowPlaying) TableName() string {
	return "channels_overlays_now_playing"
}
