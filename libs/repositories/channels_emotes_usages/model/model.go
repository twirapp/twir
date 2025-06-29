package model

import (
	"time"
)

type EmoteStatistic struct {
	EmoteName         string
	TotalUsages       uint64
	LastUsedTimestamp time.Time
}

type EmoteRange struct {
	Count     uint64
	TimeStamp time.Time
}

type EmoteUsage struct {
	ChannelID string
	UserID    string
	Emote     string
	CreatedAt time.Time
}

type EmoteUsageTopUser struct {
	ChannelID string
	UserID    string
	Count     uint64
}

type UserMostUsedEmote struct {
	Emote string
	Count uint64
}
