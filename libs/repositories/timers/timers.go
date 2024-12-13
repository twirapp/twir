package timers

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/timers/model"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (model.Timer, error)
	GetAllByChannelID(ctx context.Context, channelID string) ([]model.Timer, error)
	CountByChannelID(ctx context.Context, channelID string) (int, error)
	Create(ctx context.Context, data CreateInput) (model.Timer, error)
	UpdateByID(ctx context.Context, id string, data UpdateInput) (model.Timer, error)
	Delete(ctx context.Context, id string) error
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
	Text       string
	IsAnnounce bool
}

type UpdateInput struct {
	Name            *string
	Enabled         *bool
	TimeInterval    *int
	MessageInterval *int
	Responses       []CreateResponse
}
