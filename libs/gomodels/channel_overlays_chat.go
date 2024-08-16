package model

import (
	"time"

	"github.com/google/uuid"
)

type ChatOverlaySettings struct {
	ID                  uuid.UUID                        `gorm:"column:id;type:uuid;default:uuid_generate_v4()"        json:"id"`
	MessageHideTimeout  uint32                           `gorm:"column:message_hide_timeout" json:"message_hide_timeout"`
	MessageShowDelay    uint32                           `gorm:"column:message_show_delay" json:"message_show_delay"`
	Preset              string                           `gorm:"column:preset" json:"preset"`
	FontFamily          string                           `gorm:"column:font_family" json:"font_family"`
	FontSize            uint32                           `gorm:"column:font_size" json:"font_size"`
	FontWeight          uint32                           `gorm:"column:font_weight" json:"font_weight"`
	FontStyle           string                           `gorm:"column:font_style" json:"font_style"`
	HideCommands        bool                             `gorm:"column:hide_commands" json:"hide_commands"`
	HideBots            bool                             `gorm:"column:hide_bots" json:"hide_bots"`
	ShowBadges          bool                             `gorm:"column:show_badges" json:"show_badges"`
	ShowAnnounceBadge   bool                             `gorm:"column:show_announce_badge" json:"show_announce_badge"`
	TextShadowColor     string                           `gorm:"column:text_shadow_color" json:"text_shadow_color"`
	TextShadowSize      uint32                           `gorm:"column:text_shadow_size" json:"text_shadow_size"`
	ChatBackgroundColor string                           `gorm:"column:chat_background_color" json:"chat_background_color"`
	Direction           string                           `gorm:"column:direction" json:"direction"`
	PaddingContainer    uint32                           `gorm:"column:padding_container" json:"padding_container"`
	CreatedAt           time.Time                        `gorm:"column:created_at" json:"created_at"`
	Animation           ChatOverlaySettingsAnimationType `gorm:"column:animation" json:"animation"`

	ChannelID string    `gorm:"column:channel_id;type:text" json:"channel_id"`
	Channel   *Channels `gorm:"foreignkey:ChannelID;" json:"channel"`
}

func (c ChatOverlaySettings) TableName() string {
	return "channels_overlays_chat"
}

type ChatOverlaySettingsAnimationType string

const (
	ChatOverlaySettingsAnimationTypeDisabled ChatOverlaySettingsAnimationType = "DISABLED"
	ChatOverlaySettingsAnimationTypeDefault  ChatOverlaySettingsAnimationType = "DEFAULT"
)
