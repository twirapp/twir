package giveaways

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/twirapp/twir/libs/repositories/giveaways/model"
)

type Repository interface {
	GetManyByChannelID(ctx context.Context, channelID string) ([]model.ChannelGiveaway, error)
	GetManyActiveByChannelID(ctx context.Context, channelID string) ([]model.ChannelGiveaway, error)
	GetByChannelIDAndKeyword(
		ctx context.Context,
		channelID, keyword string,
	) (model.ChannelGiveaway, error)
	GetByID(ctx context.Context, id ulid.ULID) (model.ChannelGiveaway, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelGiveaway, error)
	Delete(ctx context.Context, id ulid.ULID) error
	Update(ctx context.Context, id ulid.ULID, input UpdateInput) (model.ChannelGiveaway, error)
}

type CreateInput struct {
	ChannelID       string
	Keyword         string
	CreatedByUserID string
}

type UpdateInput struct {
	StartedAt  *time.Time
	EndedAt    *time.Time
	Keyword    *string
	ArchivedAt *time.Time
	StoppedAt  *time.Time
}
