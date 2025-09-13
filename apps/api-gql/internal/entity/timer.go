package entity

import (
	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/bus-core/bots"
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
	ID            uuid.UUID
	Text          string
	IsAnnounce    bool
	Count         int
	AnnounceColor bots.AnnounceColor
}

var TimerNil = Timer{}
