package entity

import (
	"time"

	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type DashboardStats struct {
	StreamCategoryID   string
	StreamCategoryName string
	StreamViewers      *int
	StreamStartedAt    *time.Time
	StreamTitle        string
	StreamChatMessages int
	Followers          int
	UsedEmotes         int
	RequestedSongs     int
	Subs               int
	Platforms          []PlatformStats
}

type PlatformStats struct {
	Platform     platformentity.Platform
	IsLive       bool
	Title        *string
	CategoryID   *string
	CategoryName *string
	Viewers      *int
	Followers    *int
	StartedAt    *time.Time
	ChatMessages int
	UsedEmotes   int
	CanEditInfo  bool
}

type BotStatus struct {
	DashboardID string
	Platform    string
	ChannelName string
	IsMod       bool
	BotID       string
	BotName     string
	Enabled     bool
}
