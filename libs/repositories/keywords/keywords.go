package keywords

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/keywords/model"
)

type Repository interface {
	GetAllByChannelID(ctx context.Context, channelID uuid.UUID) ([]model.Keyword, error)
	CountByChannelID(ctx context.Context, channelID uuid.UUID) (int, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Keyword, error)
	Create(ctx context.Context, input CreateInput) (model.Keyword, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.Keyword, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	ChannelID        uuid.UUID
	Text             string
	Response         string
	Enabled          bool
	Cooldown         int
	CooldownExpireAt *time.Time
	IsReply          bool
	IsRegular        bool
	Usages           int
	RolesIDs         []uuid.UUID
	Platforms        []platform.Platform
}

type UpdateInput struct {
	Text             *string
	Response         *string
	Enabled          *bool
	Cooldown         *int
	CooldownExpireAt *time.Time
	IsReply          *bool
	IsRegular        *bool
	Usages           *int
	RolesIDs         *[]uuid.UUID
	Platforms        []platform.Platform
}
