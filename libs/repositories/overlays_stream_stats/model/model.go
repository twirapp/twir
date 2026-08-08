package model

import (
	"time"

	"github.com/google/uuid"
)

type StreamStatsOverlay struct {
	ID                   uuid.UUID `db:"id"`
	ChannelID            uuid.UUID `db:"channel_id"`
	Design               string    `db:"design"`
	Variant              string    `db:"variant"`
	ViewersEnabled       bool      `db:"viewers_enabled"`
	ViewersMode          string    `db:"viewers_mode"`
	PlatformIconsEnabled bool      `db:"platform_icons_enabled"`
	MessagesEnabled      bool      `db:"messages_enabled"`
	UptimeEnabled        bool      `db:"uptime_enabled"`
	SubscribersEnabled   bool      `db:"subscribers_enabled"`
	FollowersEnabled     bool      `db:"followers_enabled"`
	ViewersColor         string    `db:"viewers_color"`
	MessagesColor        string    `db:"messages_color"`
	UptimeColor          string    `db:"uptime_color"`
	SubscribersColor     string    `db:"subscribers_color"`
	FollowersColor       string    `db:"followers_color"`
	CounterOrder         []string  `db:"counter_order"`
	CustomHTMLEnabled    bool      `db:"custom_html_enabled"`
	CustomHTML           string    `db:"custom_html"`
	CustomCSS            string    `db:"custom_css"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`

	isNil bool
}

func (s StreamStatsOverlay) IsNil() bool {
	return s.isNil || s.ID == uuid.Nil
}

var Nil = StreamStatsOverlay{
	isNil: true,
}
