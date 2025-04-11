package model

import (
	"time"
)

type ChannelInfoHistory struct {
	ID        string
	ChannelID string
	Title     string
	Category  string
	CreatedAt time.Time
}
