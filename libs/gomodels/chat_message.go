package model

import (
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID              uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	ChannelID       string    `gorm:"column:channel_id;type:text"`
	UserID          string    `gorm:"column:user_id;type:text"`
	UserName        string    `gorm:"column:user_name;type:text"`
	UserDisplayName string    `gorm:"column:user_display_name;type:text"`
	UserColor       string    `gorm:"column:user_color;type:text"`
	Text            string    `gorm:"column:text;type:text"`
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamp"`

	User    *Users    `gorm:"foreignKey:UserID;references:ID"`
	Channel *Channels `gorm:"foreignKey:ChannelID;references:ID"`
}

func (c *ChatMessage) TableName() string {
	return "chat_messages"
}
