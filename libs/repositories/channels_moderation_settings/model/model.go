package model

import (
	"time"

	"github.com/google/uuid"
)

type ChannelModerationSettings struct {
	ID        uuid.UUID
	Type      ModerationSettingsType
	ChannelID string
	Enabled   bool
	Name      *string

	BanTime        int32
	BanMessage     string
	WarningMessage string

	CheckClips    bool
	TriggerLength int
	MaxPercentage int
	// ISO639_1
	DeniedChatLanguages []string
	ExcludedRoles       []string
	MaxWarnings         int

	DenyList                    []string
	DenyListRegexpEnabled       bool
	DenyListWordBoundaryEnabled bool
	DenyListSensitivityEnabled  bool

	OneManSpamMinimumStoredMessages int
	OneManSpamMessageMemorySeconds  int

	LanguageExcludedWords []string

	CreatedAt time.Time
	UpdatedAt time.Time

	isNil bool
}

func (c ChannelModerationSettings) IsNil() bool {
	return c.isNil
}

type ModerationSettingsType string

func (c ModerationSettingsType) String() string {
	return string(c)
}

const (
	ModerationSettingsTypeLinks       ModerationSettingsType = "links"
	ModerationSettingsTypeDenylist    ModerationSettingsType = "deny_list"
	ModerationSettingsTypeSymbols     ModerationSettingsType = "symbols"
	ModerationSettingsTypeLongMessage ModerationSettingsType = "long_message"
	ModerationSettingsTypeCaps        ModerationSettingsType = "caps"
	ModerationSettingsTypeEmotes      ModerationSettingsType = "emotes"
	ModerationSettingsTypeLanguage    ModerationSettingsType = "language"
	ModerationSettingsTypeOneManSpam  ModerationSettingsType = "one_man_spam"
)

var Nil = ChannelModerationSettings{
	isNil: true,
}
