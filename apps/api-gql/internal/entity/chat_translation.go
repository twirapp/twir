package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChatTranslation struct {
	ID                uuid.UUID
	ChannelID         string
	Enabled           bool
	TargetLanguage    string
	ExcludedLanguages []string
	UseItalic         bool
	ExcludedUsersIDs  []string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
