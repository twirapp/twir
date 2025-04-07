package model

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
	Responses                []Response `db:"responses"`
}

type Response struct {
	ID         uuid.UUID
	Text       string
	IsAnnounce bool
	TimerID    uuid.UUID
	Count      int
}

var Nil = Timer{}
