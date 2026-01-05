package giveaways

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/twirapp/twir/libs/repositories/giveaways/model"
)

type Repository interface {
	GetManyByChannelID(ctx context.Context, channelID string) ([]model.ChannelGiveaway, error)
	GetManyActiveByChannelID(ctx context.Context, channelID string) ([]model.ChannelGiveaway, error)
	GetByChannelIDAndKeyword(
		ctx context.Context,
		channelID, keyword string,
	) (model.ChannelGiveaway, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.ChannelGiveaway, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelGiveaway, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.ChannelGiveaway, error)
	UpdateStatuses(
		ctx context.Context,
		id uuid.UUID,
		input UpdateStatusInput,
	) (model.ChannelGiveaway, error)
}

type CreateInput struct {
	ChannelID       string
	Keyword         string
	CreatedByUserID string
}

type UpdateInput struct {
	StartedAt *time.Time
	Keyword   *string
	StoppedAt *time.Time
}

type UpdateStatusInput struct {
	StartedAt null.Time
	StoppedAt null.Time
}
