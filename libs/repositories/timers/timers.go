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
	Count(ctx context.Context, input CountInput) (int64, error)
	GetMany(ctx context.Context, input GetManyInput) ([]model.Timer, error)
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

type CountInput struct {
	ChannelID *string
	Enabled   *bool
}

type GetManyInput struct {
	ChannelID *string
	Enabled   *bool
	Limit     int
	Offset    int
}
