package entity

type EmoteStatistic struct {
	EmoteName         string
	TotalUsages       uint64
	LastUsedTimestamp int64
}

type EmoteRange struct {
	Count     uint64
	TimeStamp int64
}
