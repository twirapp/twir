package model

import (
	"github.com/google/uuid"
)

type ChannelGamesSeppuku struct {
	ID                uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid()"`
	ChannelID         string    `gorm:"column:channel_id;type:text"`
	Enabled           bool      `gorm:"column:enabled;type:bool"`
	TimeoutSeconds    int8      `gorm:"column:timeout_seconds;type:int2"`
	TimeoutModerators bool      `gorm:"column:timeout_moderators;type:bool"`
	Message           string    `gorm:"column:message;type:text"`
	MessageModerators string    `gorm:"column:message_moderators;type:text"`
}

func (ChannelGamesSeppuku) TableName() string {
	return "channels_games_seppuku"
}
