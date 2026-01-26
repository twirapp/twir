package timers

import (
	"context"

	"github.com/google/uuid"
	timecrsentity "github.com/twirapp/twir/libs/entities/timers"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (timecrsentity.Timer, error)
	GetAllByChannelID(ctx context.Context, channelID string) ([]timecrsentity.Timer, error)
	CountByChannelID(ctx context.Context, channelID string) (int, error)
	Create(ctx context.Context, data CreateInput) (timecrsentity.Timer, error)
	UpdateByID(ctx context.Context, id uuid.UUID, data UpdateInput) (timecrsentity.Timer, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context, input CountInput) (int64, error)
	GetMany(ctx context.Context, input GetManyInput) ([]timecrsentity.Timer, error)
}

type CreateInput struct {
	ChannelID       string
	Name            string
	Enabled         bool
	OfflineEnabled  bool
	TimeInterval    int
	MessageInterval int
	Responses       []CreateResponse
}

type CreateResponse struct {
	Text          string
	IsAnnounce    bool
	Count         int
	AnnounceColor timecrsentity.AnnounceColor
}

type UpdateInput struct {
	Name            *string
	Enabled         *bool
	OfflineEnabled  *bool
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
