package greetings

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/greetings/model"
)

type Repository interface {
	GetManyByChannelID(ctx context.Context, channelID string) ([]model.Greeting, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Greeting, error)
	Create(ctx context.Context, input CreateInput) (model.Greeting, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.Greeting, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	ChannelID string
	UserID    string
	Enabled   bool
	Text      string
	IsReply   bool
	Processed bool
}

type UpdateInput struct {
	UserID    *string
	Enabled   *bool
	Text      *string
	IsReply   *bool
	Processed *bool
}
