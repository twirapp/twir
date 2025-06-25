package entity

import (
	"time"
)

type EmoteStatistic struct {
	EmoteName         string
	TotalUsages       uint64
	LastUsedTimestamp int64
}

type EmoteRange struct {
	Count     uint64
	TimeStamp int64
}

type EmoteStatisticTopUser struct {
	UserID string
	Count  int
}

type EmoteStatisticUserUsage struct {
	UserID string
	Date   time.Time
}
