package timers

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/timers/model"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Timer, error)
	GetAllByChannelID(ctx context.Context, channelID string) ([]model.Timer, error)
	CountByChannelID(ctx context.Context, channelID string) (int, error)
	Create(ctx context.Context, data CreateInput) (model.Timer, error)
	UpdateByID(ctx context.Context, id uuid.UUID, data UpdateInput) (model.Timer, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	ChannelID       string
	Name            string
	Enabled         bool
	TimeInterval    int
	MessageInterval int
	Responses       []CreateResponse
}

type CreateResponse struct {
	Text          string
	IsAnnounce    bool
	Count         int
	AnnounceColor model.AnnounceColor
}

type UpdateInput struct {
	Name            *string
	Enabled         *bool
	TimeInterval    *int
	MessageInterval *int
	Responses       []CreateResponse
}
