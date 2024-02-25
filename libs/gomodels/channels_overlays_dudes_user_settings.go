package model

import (
	"time"

	"github.com/google/uuid"
)

type ChannelsOverlaysDudesUserSettings struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;column:id"  json:"id"`
	ChannelID  string    `gorm:"type:text;not null;column:channel_id" json:"channelId"`
	UserID     string    `gorm:"type:text;not null;column:user_id" json:"userId"`
	DudeColor  *string   `gorm:"type:text;column:dude_color" json:"dudeColor"`
	DudeSprite *string   `gorm:"type:text;column:dude_sprite" json:"dudeSprite"`

	CreatedAt time.Time `gorm:"type:datetime;not null;" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;" json:"updatedAt"`

	Channel *Channels `gorm:"foreignKey:ChannelID" json:"channel"`
	User    *Users    `gorm:"foreignKey:UserID" json:"user"`
}

func (ChannelsOverlaysDudesUserSettings) TableName() string {
	return "channels_overlays_dudes_user_settings"
}
