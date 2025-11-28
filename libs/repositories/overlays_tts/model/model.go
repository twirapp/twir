package model

import (
	"time"

	"github.com/google/uuid"
)

type TTSOverlay struct {
	ID        uuid.UUID
	ChannelID string
	CreatedAt time.Time
	UpdatedAt time.Time
	Settings  *TTSOverlaySettings
}

type TTSOverlaySettings struct {
	Enabled                            bool     `json:"enabled"`
	Voice                              string   `json:"voice"`
	DisallowedVoices                   []string `json:"disallowed_voices"`
	Pitch                              int32    `json:"pitch"`
	Rate                               int32    `json:"rate"`
	Volume                             int32    `json:"volume"`
	DoNotReadTwitchEmotes              bool     `json:"do_not_read_twitch_emotes"`
	DoNotReadEmoji                     bool     `json:"do_not_read_emoji"`
	DoNotReadLinks                     bool     `json:"do_not_read_links"`
	AllowUsersChooseVoiceInMainCommand bool     `json:"allow_users_choose_voice_in_main_command"`
	MaxSymbols                         int32    `json:"max_symbols"`
	ReadChatMessages                   bool     `json:"read_chat_messages"`
	ReadChatMessagesNicknames          bool     `json:"read_chat_messages_nicknames"`
}

type TTSUserSettings struct {
	ID        uuid.UUID
	ChannelID string
	UserID    string
	Voice     string
	Rate      int32
	Pitch     int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

var Nil = TTSOverlay{}
var NilUserSettings = TTSUserSettings{}
