package entity

import (
	"github.com/google/uuid"
)

type Timer struct {
	ID                       uuid.UUID
	ChannelID                string
	Name                     string
	Enabled                  bool
	TimeInterval             int
	MessageInterval          int
	LastTriggerMessageNumber int
	Responses                []Response
}

type Response struct {
	ID         uuid.UUID
	Text       string
	IsAnnounce bool
}

var TimerNil = Timer{}
