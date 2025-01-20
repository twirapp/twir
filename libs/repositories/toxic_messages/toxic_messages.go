package toxic_messages

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
}

type CreateInput struct {
	ChannelID     string
	ReplyToUserID *string
	Text          string
}
