package model

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ChannelGames8Ball struct {
	ID        uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid()" json:"id"`
	ChannelId string         `gorm:"column:channel_id;type:text" json:"channelId"`
	Enabled   bool           `gorm:"column:enabled" json:"enabled"`
	Answers   pq.StringArray `gorm:"column:answers;type:text[]" json:"answers"`
}

func (ChannelGames8Ball) TableName() string {
	return "channels_games_8ball"
}
