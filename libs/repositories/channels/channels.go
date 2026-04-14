package channels

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/channels/model"
)

type Repository interface {
	GetMany(ctx context.Context, input GetManyInput) ([]model.Channel, error)
	GetByID(ctx context.Context, channelID string) (model.Channel, error)
	GetByUserIDAndPlatform(ctx context.Context, userID uuid.UUID, platformVal platform.Platform) (model.Channel, error)
	GetByPlatformUserID(ctx context.Context, plat platform.Platform, platformUserID string) (model.Channel, error)
	GetCount(ctx context.Context, input GetCountInput) (int, error)
	Update(ctx context.Context, channelID string, input UpdateInput) (model.Channel, error)
	Create(ctx context.Context, input CreateInput) (model.Channel, error)
}

type CreateInput struct {
	UserID   uuid.UUID
	BotID    string
	Platform platform.Platform
}

type UpdateInput struct {
	IsEnabled *bool
	IsBotMod  *bool
}

type GetManyInput struct {
	Enabled *bool
	PerPage int
	Page    int
}

type GetCountInput struct {
	OnlyEnabled bool
}
