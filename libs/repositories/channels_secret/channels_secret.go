package channels_secret

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/channels_secret/model"
)

var ErrNotFound = errors.New("secret not found")

type Repository interface {
	GetAllByChannelID(ctx context.Context, channelID string) ([]model.ChannelSecret, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.ChannelSecret, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelSecret, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.ChannelSecret, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	ChannelID   string
	Name        string
	Description *string
	Value       string
}

type UpdateInput struct {
	Name        *string
	Description *string
	Value       *string
}
