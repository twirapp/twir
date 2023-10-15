package model

import (
	"time"
)

type DiscordSendedNotification struct {
	ID               string `gorm:"primaryKey;column:id;type:TEXT;default:uuid_generate_v4()" json:"id"`
	GuildID          string `gorm:"column:guild_id;type:TEXT;"      json:"guildId"`
	MessageID        string `gorm:"column:message_id;type:TEXT;"    json:"messageId"`
	TwitchChannelID  string `gorm:"column:channel_id;type:TEXT;"    json:"channelId"`
	DiscordChannelID string `gorm:"column:discord_channel_id;type:TEXT;"    json:"discordChannelId"`

	CreatedAt time.Time `gorm:"column:created_at;data:timestamp;default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;data:timestamp;default:now()" json:"updatedAt"`
}

func (DiscordSendedNotification) TableName() string {
	return "discord_sended_notifications"
}
