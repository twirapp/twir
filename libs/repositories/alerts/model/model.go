package model

import (
	"github.com/google/uuid"
)

type Alert struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ChannelID    string    `json:"channel_id"`
	AudioID      *string   `json:"audio_id"`
	AudioVolume  int       `json:"audio_volume"`
	CommandIDS   []string  `json:"command_ids"`
	RewardIDS    []string  `json:"reward_ids"`
	GreetingsIDS []string  `json:"greetings_ids"`
	KeywordsIDS  []string  `json:"keywords_ids"`
}

var Nil = Alert{}
