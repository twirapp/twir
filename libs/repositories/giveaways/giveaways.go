package giveaways

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	channels_giveaways "github.com/twirapp/twir/libs/entities/channels_giveaways"
)

type Repository interface {
	GetManyByChannelID(ctx context.Context, channelID string) ([]channels_giveaways.Giveaway, error)
	GetManyActiveByChannelID(ctx context.Context, channelID string) ([]channels_giveaways.Giveaway, error)
	GetByChannelIDAndKeyword(
		ctx context.Context,
		channelID, keyword string,
	) (channels_giveaways.Giveaway, error)
	GetByID(ctx context.Context, id uuid.UUID) (channels_giveaways.Giveaway, error)
	Create(ctx context.Context, input CreateInput) (channels_giveaways.Giveaway, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (channels_giveaways.Giveaway, error)
	UpdateStatuses(
		ctx context.Context,
		id uuid.UUID,
		input UpdateStatusInput,
	) (channels_giveaways.Giveaway, error)
}

type CreateInput struct {
	ChannelID            string
	Type                 channels_giveaways.GiveawayType
	Keyword              *string
	MinWatchedTime       *int64
	MinMessages          *int32
	MinUsedChannelPoints *int64
	MinFollowDuration    *int64
	RequireSubscription  bool
	CreatedByUserID      string
}

type UpdateInput struct {
	StartedAt            *time.Time
	Keyword              *string
	StoppedAt            *time.Time
	MinWatchedTime       *int64
	MinMessages          *int32
	MinUsedChannelPoints *int64
	MinFollowDuration    *int64
	RequireSubscription  *bool
}

type UpdateStatusInput struct {
	StartedAt null.Time
	StoppedAt null.Time
}
