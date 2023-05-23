package model

import "time"

type ChannelEmoteUsage struct {
	ID        string    `gorm:"primary_key;column:id;type:TEXT;" json:"id"`
	ChannelID string    `gorm:"column:channelId;type:text"              json:"channelId"`
	UserID    string    `gorm:"column:userId;type:text"                 json:"userId"`
	Emote     string    `gorm:"column:emote;type:text"                 json:"emote"`
	CreatedAt time.Time `gorm:"column:createdAt;type:timestamp"         json:"createdAt"`

	Channel *Channels `gorm:"foreignKey:ID" json:"channel"`
	User    *Users    `gorm:"foreignKey:ID" json:"user"`
}

type ChannelEmoteUsageWithCount struct {
	*ChannelEmoteUsage

	Count int `gorm:"count"`
}

func (c *ChannelEmoteUsage) TableName() string {
	return "channels_emotes_usages"
}
