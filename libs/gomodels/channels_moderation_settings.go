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
	ModerationSettingsTypeDenylist                           = "denylist"
	ModerationSettingsTypeSymbols                            = "symbols"
	ModerationSettingsTypeLongMessage                        = "longMessage"
	ModerationSettingsTypeCaps                               = "caps"
	ModerationSettingsTypeEmotes                             = "emotes"
	ModerationSettingsTypeLanguage                           = "language"
)

func (c ModerationSettingsType) String() string {
	return string(c)
}

type ChannelsModerationSettings struct {
	ID                    string                 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;"`
	Type                  ModerationSettingsType `gorm:"column:type;type:VARCHAR;"`
	ChannelID             string                 `gorm:"column:channel_id;type:TEXT;"`
	Enabled               bool                   `gorm:"column:enabled;type:BOOL;default:false;"`
	BanTime               int32                  `gorm:"column:ban_time;type:INT4;default:600;"`
	BanMessage            string                 `gorm:"column:ban_message;type:TEXT;"`
	WarningMessage        string                 `gorm:"column:warning_message;type:TEXT;"`
	CheckClips            bool                   `gorm:"column:check_clips;type:BOOL;default:false;"`
	TriggerLength         null.Int               `gorm:"column:trigger_length;type:INT4;default:300;"`
	MaxPercentage         null.Int               `gorm:"column:max_percentage;type:INT4;default:50;"`
	DenyList              pq.StringArray         `gorm:"column:deny_list;type:JSONB;default:[];"`
	AcceptedChatLanguages pq.StringArray         `gorm:"column:accepted_chat_languages;type:JSONB;default:[];"`
	ExcludedRoles         pq.StringArray         `gorm:"column:excluded_roles;type:JSONB;default:[];"`
}

func (c *ChannelsModerationSettings) TableName() string {
	return "channels_moderation_settings"
}
