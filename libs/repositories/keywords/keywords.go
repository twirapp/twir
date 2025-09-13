package keywords

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/keywords/model"
)

type Repository interface {
	GetAllByChannelID(ctx context.Context, channelID string) ([]model.Keyword, error)
	CountByChannelID(ctx context.Context, channelID string) (int, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Keyword, error)
	Create(ctx context.Context, input CreateInput) (model.Keyword, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.Keyword, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	ChannelID        string
	Text             string
	Response         string
	Enabled          bool
	Cooldown         int
	CooldownExpireAt *time.Time
	IsReply          bool
	IsRegular        bool
	Usages           int
	RolesIDs         []uuid.UUID
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
}
