package channels

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/channels/model"
)

type Repository interface {
	GetMany(ctx context.Context, input GetManyInput) ([]model.Channel, error)
	GetByID(ctx context.Context, channelID uuid.UUID) (model.Channel, error)
	GetByApiKey(ctx context.Context, apiKey string) (model.Channel, error)
	GetByTwitchUserID(ctx context.Context, twitchUserID uuid.UUID) (model.Channel, error)
	GetByTwitchPlatformID(ctx context.Context, twitchPlatformID string) (model.Channel, error)
	GetByKickUserID(ctx context.Context, kickUserID uuid.UUID) (model.Channel, error)
	GetByKickPlatformID(ctx context.Context, kickPlatformID string) (model.Channel, error)
	GetBySlug(ctx context.Context, opts GetBySlugInput) (model.Channel, error)
	GetCount(ctx context.Context, input GetCountInput) (int, error)
	Update(ctx context.Context, channelID uuid.UUID, input UpdateInput) (model.Channel, error)
	Create(ctx context.Context, input CreateInput) (model.Channel, error)
}

type CreateInput struct {
	TwitchUserID     *uuid.UUID
	KickUserID       *uuid.UUID
	TwitchBotEnabled bool
	KickBotEnabled   bool
	BotID            string
	KickBotID        *uuid.UUID
}

type UpdateInput struct {
	IsEnabled        *bool
	IsBotMod         *bool
	TwitchUserID     *uuid.UUID
	KickUserID       *uuid.UUID
	TwitchBotEnabled *bool
	KickBotEnabled   *bool
	KickBotID        *uuid.UUID
}

type GetManyInput struct {
	Enabled          *bool
	TwitchBotEnabled *bool
	KickBotEnabled   *bool
	AnyBotEnabled    *bool
	HasKickUserID    *bool
	HasTwitchUserID  *bool
	PerPage          int
	Page             int
}

type GetCountInput struct {
	OnlyEnabled       bool
	OnlyTwitchEnabled bool
}

type GetBySlugInput struct {
	Slug     string
	Platform *platform.Platform
}
