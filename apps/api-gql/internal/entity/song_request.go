package entity

import (
	"time"
)

type SongRequestPublic struct {
	Title           string
	UserID          string
	CreatedAt       time.Time
	SongLink        string
	DurationSeconds int
}

type SongRequestPlaybackState struct {
	VideoID   string
	Title     string
	Position  float64
	IsPlaying bool
	Volume    int
	UpdatedAt int64
}

type SongRequestQueueItem struct {
	ID                   string
	Title                string
	SongLink             string
	DurationSeconds      int
	OrderedByName        string
	OrderedByDisplayName string
	QueuePosition        int
	CreatedAt            string
}

type SongRequestWidgetData struct {
	PlaybackState *SongRequestPlaybackState
	Queue         []SongRequestQueueItem
	Volume        int
}

type ChannelByApiKeyResult struct {
	ID           string
	TwitchUserID *string
	KickUserID   *string
}
