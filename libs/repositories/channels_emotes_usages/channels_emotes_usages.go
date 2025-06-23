package channels_emotes_usages

import (
	"context"
)

type Repository interface {
	CreateMany(ctx context.Context, inputs []ChannelEmoteUsageInput) error
	Count(ctx context.Context, input CountInput) (uint64, error)
}

type ChannelEmoteUsageInput struct {
	ChannelID string
	UserID    string
	Emote     string
}

type CountInput struct {
	ChannelID *string
	UserID    *string
}
