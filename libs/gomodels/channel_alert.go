package model

import "github.com/guregu/null"

type ChannelAlert struct {
	ID          string      `gorm:"primary_key;column:id;type:TEXT;" json:"id"`
	ChannelID   string      `gorm:"column:channel_id;type:text"              json:"channel_id"`
	Name        string      `gorm:"column:name;type:text"              json:"name"`
	AudioID     null.String `gorm:"column:audio_id;type:text"              json:"audio_id"`
	AudioVolume int         `gorm:"column:audio_volume;type:int2;default:100"  json:"audio_volume"`

	Channel *Channels    `gorm:"foreignKey:ChannelID" json:"channel"`
	Audio   *ChannelFile `gorm:"foreignKey:AudioID" json:"audio"`
}

func (C ChannelAlert) TableName() string {
	return "channels_alerts"
}
