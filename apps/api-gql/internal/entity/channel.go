package entity

import "github.com/google/uuid"

type Channel struct {
	ID             uuid.UUID
	IsEnabled      bool
	IsTwitchBanned bool
	IsBotMod       bool
	BotID          string
}

var ChannelNil = Channel{}
