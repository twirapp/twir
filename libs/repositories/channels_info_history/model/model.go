package model

import (
	"time"

	"github.com/twirapp/twir/libs/entities/platform"
)

type ChannelInfoHistory struct {
	ID        string
	ChannelID string
	Platform  platform.Platform
	Title     string
	Category  string
	CreatedAt time.Time
}
