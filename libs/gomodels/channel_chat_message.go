package model

import "time"

type ChannelChatMessage struct {
	MessageId    string    `gorm:"primary_key;column:messageId;type:TEXT;" json:"messageId"`
	ChannelId    string    `gorm:"column:channelId;type:text"              json:"channelId"`
	UserId       string    `gorm:"column:userId;type:text"                 json:"userId"`
	UserName     string    `gorm:"column:userName;type:text"               json:"userName"`
	Text         string    `gorm:"column:text;type:text"                   json:"text"`
	CanBeDeleted bool      `gorm:"column:canBeDeleted;type:bool"           json:"canBeDeleted"`
	CreatedAt    time.Time `gorm:"column:createdAt;type:timestamp"         json:"createdAt"`

	Channel *Channels `gorm:"foreignKey:ID" json:"channel"`
	User    *Users    `gorm:"foreignKey:ID" json:"user"`
}

func (c *ChannelChatMessage) TableName() string {
	return "channels_messages"
}
