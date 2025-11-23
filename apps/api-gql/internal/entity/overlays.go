package entity

import (
	"time"

	"github.com/google/uuid"
)

type TTSUserSettings struct {
	UserID         string
	Rate           int
	Pitch          int
	Voice          string
	IsChannelOwner bool
}

type BeRightBackOverlay struct {
	ID        uuid.UUID
	ChannelID string
	CreatedAt time.Time
	UpdatedAt time.Time
	Settings  BeRightBackOverlaySettings
}

type BeRightBackOverlaySettings struct {
	Text            string
	Late            BeRightBackOverlayLateSettings
	BackgroundColor string
	FontSize        int32
	FontColor       string
	FontFamily      string
}

type BeRightBackOverlayLateSettings struct {
	Enabled        bool
	Text           string
	DisplayBrbTime bool
}

type TTSOverlay struct {
	ID        uuid.UUID
	ChannelID string
	CreatedAt time.Time
	UpdatedAt time.Time
	Settings  TTSOverlaySettings
}

type TTSOverlaySettings struct {
	Enabled                            bool
	Voice                              string
	DisallowedVoices                   []string
	Pitch                              int32
	Rate                               int32
	Volume                             int32
	DoNotReadTwitchEmotes              bool
	DoNotReadEmoji                     bool
	DoNotReadLinks                     bool
	AllowUsersChooseVoiceInMainCommand bool
	MaxSymbols                         int32
	ReadChatMessages                   bool
	ReadChatMessagesNicknames          bool
}

