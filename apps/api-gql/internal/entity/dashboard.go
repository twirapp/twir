package entity

import (
	"time"
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
