package chat_messages

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/chat_messages/model"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
	CreateMany(ctx context.Context, input []CreateInput) error
	GetMany(ctx context.Context, input GetManyInput) ([]model.ChatMessage, error)
}

type CreateInput struct {
	ID              string
	ChannelID       string
	UserID          string
	Text            string
	UserName        string
	UserDisplayName string
	UserColor       string
}

type GetManyInput struct {
	Page    int
	PerPage int

	ChannelID    *string
	UserNameLike *string
	TextLike     *string
}
