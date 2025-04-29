package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
	DeniedChatLanguages pq.StringArray
	ExcludedRoles       pq.StringArray
	MaxWarnings         int

	DenyList                    pq.StringArray
	DenyListRegexpEnabled       bool
	DenyListWordBoundaryEnabled bool
	DenyListSensitivityEnabled  bool

	CreatedAt time.Time
	UpdatedAt time.Time
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
)
