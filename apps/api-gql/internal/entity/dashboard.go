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
	IsMod   bool
	BotID   string
	BotName string
	Enabled bool
}
