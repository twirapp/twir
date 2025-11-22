package model

import (
	"time"

	"github.com/google/uuid"
)

type TTS struct {
	ID                                 uuid.UUID
	ChannelID                          string
	UserID                             *string
	CreatedAt                          time.Time
	UpdatedAt                          time.Time
	Enabled                            *bool
	Rate                               int
	Volume                             int
	Pitch                              int
	Voice                              string
	AllowUsersChooseVoiceInMainCommand bool
	MaxSymbols                         int
	DisallowedVoices                   []string
	DoNotReadEmoji                     bool
	DoNotReadTwitchEmotes              bool
	DoNotReadLinks                     bool
	ReadChatMessages                   bool
	ReadChatMessagesNicknames          bool
}

var Nil = TTS{}
