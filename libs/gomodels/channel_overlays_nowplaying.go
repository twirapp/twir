package model

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/twirapp/twir/libs/types/types/api/overlays"
)

type ChannelOverlayNowPlaying struct {
	ID              uuid.UUID                               `gorm:"primary_key;column:id;type:UUID;"  json:"id"`
	Preset          overlays.ChannelOverlayNowPlayingPreset `gorm:"column:preset;type:varchar;"  json:"preset"`
	FontFamily      string                                  `gorm:"column:font_family;type:text" json:"fontFamily"`
	FontWeight      uint32                                  `gorm:"column:font_weight;type:int" json:"fontWeight"`
	BackgroundColor string                                  `gorm:"column:background_color;type:text" json:"backgroundColor"`
	ShowImage       bool                                    `gorm:"column:show_image;type:boolean" json:"showImage"`
	HideTimeout     null.Int                                `gorm:"column:hide_timeout;type:int" json:"hideTimeout"`

	ChannelID string    `gorm:"column:channel_id;type:text" json:"channelId"`
	Channel   *Channels `gorm:"foreignkey:ChannelID;" json:"channel"`
}

func (c ChannelOverlayNowPlaying) TableName() string {
	return "channels_overlays_now_playing"
}
