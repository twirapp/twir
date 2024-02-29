package model

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type ChannelPublicSettings struct {
	ID          uuid.UUID                         `gorm:"column:id;type:uuid;default:uuid_generate_v4()"`
	ChannelID   string                            `gorm:"column:channel_id;type:TEXT;"`
	Description null.String                       `gorm:"column:description;type:varchar(1000)"`
	SocialLinks []ChannelPublicSettingsSocialLink `gorm:"foreignKey:SettingsID"`
}

func (ChannelPublicSettings) TableName() string {
	return "channels_public_settings"
}

type ChannelPublicSettingsSocialLink struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4()"`
	Title      string    `gorm:"column:title;type:varchar(30)"`
	Href       string    `gorm:"column:href;type:varchar(500)"`
	SettingsID uuid.UUID `gorm:"column:settings_id;type:uuid"`
}

func (ChannelPublicSettingsSocialLink) TableName() string {
	return "channels_public_settings_links"
}
