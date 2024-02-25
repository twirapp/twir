package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type ChannelsOverlaysDudesUserSettings struct {
	ID        uuid.UUID   `gorm:"type:uuid;primary_key;"  json:"id"`
	ChannelID string      `gorm:"type:text;not null;" json:"channelId"`
	UserID    string      `gorm:"type:text;not null;" json:"userId"`
	DudeColor null.String `gorm:"type:text;not null;" json:"dudeColor"`

	CreatedAt time.Time `gorm:"type:datetime;not null;" json:"createdAt"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;" json:"updatedAt"`

	Channel *Channels `gorm:"foreignKey:ChannelID" json:"channel"`
	User    *Users    `gorm:"foreignKey:UserID" json:"user"`
}

func (ChannelsOverlaysDudesUserSettings) TableName() string {
	return "channels_overlays_dudes_user_settings"
}
