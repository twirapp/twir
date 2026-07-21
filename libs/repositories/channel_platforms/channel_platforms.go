package channel_platforms

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/channel_platforms/model"
)

var ErrNotFound = errors.New("channel platform binding not found")

type Repository interface {
	Create(ctx context.Context, input CreateInput) (model.ChannelPlatform, error)
	GetByChannelAndPlatform(ctx context.Context, channelID uuid.UUID, platform platform.Platform) (model.ChannelPlatform, error)
	GetByPlatformChannelID(ctx context.Context, platform platform.Platform, platformChannelID string) (model.ChannelPlatform, error)
	ListByChannelID(ctx context.Context, channelID uuid.UUID) ([]model.ChannelPlatform, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.ChannelPlatform, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	ChannelID         uuid.UUID
	Platform          platform.Platform
	UserID            uuid.UUID
	PlatformChannelID string
	Enabled           bool
	BotUserID         *uuid.UUID
	BotConfig         json.RawMessage
}

type UpdateInput struct {
	UserID            uuid.UUID
	PlatformChannelID string
	Enabled           bool
	BotUserID         *uuid.UUID
	BotConfig         json.RawMessage
}
