package sentmessages

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
}

type CreateInput struct {
	MessageTwitchID string
	Content         string
	ChannelID       string
	SenderID        string
}
