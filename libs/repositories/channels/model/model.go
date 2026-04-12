package model

import (
	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
)

type Channel struct {
	ID             string
	Platform       platform.Platform
	UserID         uuid.UUID
	IsEnabled      bool
	IsTwitchBanned bool
	IsBotMod       bool
	BotID          string

	isNil bool
}

func (c Channel) IsNil() bool {
	return c.isNil
}

var Nil = Channel{
	isNil: true,
}
