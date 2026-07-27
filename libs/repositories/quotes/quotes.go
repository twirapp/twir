package quotes

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/quotes/model"
)

type Repository interface {
	GetAllByChannelID(ctx context.Context, channelID uuid.UUID) ([]model.Quote, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Quote, error)
	Create(ctx context.Context, input CreateInput) (model.Quote, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.Quote, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	ChannelID   uuid.UUID
	Text        string
	CreatorID   *string
	CreatorName *string
	GameID      *string
	GameName    *string
}

type UpdateInput struct {
	Text *string
}
