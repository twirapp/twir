package model

import "time"

type ChannelChatMessage struct {
	ID           uint      `gorm:"primary_key;column:id;type:serial"         json:"id"`
	MessageId    string    `gorm:"primary_key;column:message_id;type:TEXT;" json:"messageId"`
	ChannelId    string    `gorm:"column:channel_id;type:text"              json:"channelId"`
	UserId       string    `gorm:"column:user_id;type:text"                 json:"userId"`
	UserName     string    `gorm:"column:user_name;type:text"               json:"userName"`
	Text         string    `gorm:"column:text;type:text"                   json:"text"`
	CanBeDeleted bool      `gorm:"column:can_be_deleted;type:bool"           json:"canBeDeleted"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp"         json:"createdAt"`

	Channel *Channels `gorm:"foreignKey:ID" json:"channel"`
	User    *Users    `gorm:"foreignKey:ID" json:"user"`
}

func (c *ChannelChatMessage) TableName() string {
	return "channels_messages"
}
