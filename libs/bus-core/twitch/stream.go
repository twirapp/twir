package twitch

import (
	"time"
)

type StreamUpdateMessage struct {
	ChannelID string `json:"channelId"`
	Title     string `json:"title"`
	Category  string `json:"category"`
}

type StreamOnlineMessage struct {
	ChannelID string `json:"channelId"`
	StreamID  string `json:"streamId"`

	CategoryName string    `json:"categoryName"`
	CategoryID   string    `json:"categoryId"`
	Title        string    `json:"title"`
	Viewers      int       `json:"viewers"`
	StartedAt    time.Time `json:"startedAt"`
}

type StreamOfflineMessage struct {
	ChannelID string    `json:"channelId"`
	StartedAt time.Time `json:"startedAt"`
}
