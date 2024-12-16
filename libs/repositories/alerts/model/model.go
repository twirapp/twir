package model

import (
	"github.com/google/uuid"
)

type Alert struct {
	ID           uuid.UUID
	Name         string
	ChannelID    string
	AudioID      *string
	AudioVolume  int
	CommandIDS   []string
	RewardIDS    []string
	GreetingsIDS []string
	KeywordsIDS  []string
}

var Nil = Alert{}
