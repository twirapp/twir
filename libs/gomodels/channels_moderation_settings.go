package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type ModerationSettingsType string

const (
	ModerationSettingsTypeLinks       ModerationSettingsType = "links"
	ModerationSettingsTypeDenylist                           = "deny_list"
	ModerationSettingsTypeSymbols                            = "symbols"
	ModerationSettingsTypeLongMessage                        = "long_message"
	ModerationSettingsTypeCaps                               = "caps"
	ModerationSettingsTypeEmotes                             = "emotes"
	ModerationSettingsTypeLanguage                           = "language"
)

func (c ModerationSettingsType) String() string {
	return string(c)
}

type ChannelModerationSettings struct {
	ID        string                 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	Type      ModerationSettingsType `gorm:"column:type;type:VARCHAR;" json:"type"`
	ChannelID string                 `gorm:"column:channel_id;type:TEXT;" json:"channel_id"`
	Enabled   bool                   `gorm:"column:enabled;type:BOOL;default:false;" json:"enabled"`

	BanTime        int32  `gorm:"column:ban_time;type:INT4;default:600;" json:"ban_time"`
	BanMessage     string `gorm:"column:ban_message;type:TEXT;" json:"ban_message"`
	WarningMessage string `gorm:"column:warning_message;type:TEXT;" json:"warning_message"`

	CheckClips    bool           `gorm:"column:check_clips;type:BOOL;default:false;" json:"check_clips"`
	TriggerLength int            `gorm:"column:trigger_length;type:INT4;default:300;" json:"trigger_length"`
	MaxPercentage int            `gorm:"column:max_percentage;type:INT4;default:50;" json:"max_percentage"`
	DenyList      pq.StringArray `gorm:"column:deny_list;type:JSONB;default:[];" json:"deny_list"`
	// ISO639_1
	DeniedChatLanguages pq.StringArray `gorm:"column:denied_chat_languages;type:JSONB;default:[];" json:"accepted_chat_languages"`
	ExcludedRoles       pq.StringArray `gorm:"column:excluded_roles;type:JSONB;default:[];" json:"excluded_roles"`
	MaxWarnings         int            `gorm:"column:max_warnings;type:INT4;default:0;" json:"max_warnings"`

	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMPTZ;default:now();" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:TIMESTAMPTZ;default:now();" json:"updated_at"`
}

func (c *ChannelModerationSettings) TableName() string {
	return "channels_moderation_settings"
}
