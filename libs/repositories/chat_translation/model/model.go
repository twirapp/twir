package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type ChatTranslation struct {
	ID                ulid.ULID
	ChannelID         string
	Enabled           bool
	TargetLanguage    string
	ExcludedLanguages []string
	UseItalic         bool
	ExcludedUsersIDs  []string

	CreatedAt time.Time
	UpdatedAt time.Time
}

var ChatTranslationNil = ChatTranslation{}
