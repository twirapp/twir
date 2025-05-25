package channels_emotes_usages

import (
	"context"
	"time"
)

type Repository interface {
	CreateMany(ctx context.Context, inputs []ChannelEmoteUsageInput) error
}

type ChannelEmoteUsageInput struct {
	ID        string
	ChannelID string
	UserID    string
	Emote     string
	CreatedAt time.Time
}
