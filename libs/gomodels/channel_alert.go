package model

import (
	"github.com/guregu/null"
	"github.com/lib/pq"
)

type ChannelAlert struct {
	ID           string         `gorm:"primary_key;column:id;type:TEXT;" json:"id"`
	ChannelID    string         `gorm:"column:channel_id;type:text"              json:"channel_id"`
	Name         string         `gorm:"column:name;type:text"              json:"name"`
	AudioID      null.String    `gorm:"column:audio_id;type:text"              json:"audio_id"`
	AudioVolume  int            `gorm:"column:audio_volume;type:int2;default:100"  json:"audio_volume"`
	CommandIDS   pq.StringArray `gorm:"column:command_ids;type:text[];default:[]"`
	RewardIDS    pq.StringArray `gorm:"column:reward_ids;type:text[];default:[]"`
	GreetingsIDS pq.StringArray `gorm:"column:greetings_ids;type:text[];default:[]"`
	KeywordsIDS  pq.StringArray `gorm:"column:keywords_ids;type:text[];default:[]"`

	Channel *Channels    `gorm:"foreignKey:ChannelID" json:"channel"`
	Audio   *ChannelFile `gorm:"foreignKey:AudioID" json:"audio"`
}

func (C ChannelAlert) TableName() string {
	return "channels_alerts"
}
