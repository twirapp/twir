package model

import (
	"github.com/google/uuid"
)

type ChannelGamesRussianRoulette struct {
	ID                    uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid()"`
	ChannelID             string    `gorm:"column:channel_id;type:text"`
	Enabled               bool      `gorm:"column:enabled;type:bool" json:"enabled"`
	CanBeUsedByModerators bool      `gorm:"column:can_be_used_by_moderators;type:bool" json:"canBeUsedByModerator"`
	TimeoutSeconds        int       `gorm:"column:timeout_seconds;type:int" json:"timeoutSeconds"`
	DecisionSeconds       int       `gorm:"column:decision_seconds;type:int" json:"decisionSeconds"`
	TumberSize            int       `gorm:"column:tumber_size;type:int" json:"tumberSize"`
	ChargedBullets        int       `gorm:"column:charged_bullets;type:int" json:"chargedBullets"`

	InitMessage    string `gorm:"column:init_message;type:text" json:"initMessage"`
	SurviveMessage string `gorm:"column:survive_message;type:text" json:"surviveMessage"`
	DeathMessage   string `gorm:"column:death_message;type:text" json:"deathMessage"`
}

func (ChannelGamesRussianRoulette) TableName() string {
	return "channels_games_russian_roulette"
}
