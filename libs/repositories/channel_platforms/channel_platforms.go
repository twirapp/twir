package channel_platforms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/channel_platforms/model"
)

var (
	ErrNotFound              = errors.New("channel platform binding not found")
	ErrInvalidBotConfigPatch = errors.New("channel platform bot config patch must be a JSON object")
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) (model.ChannelPlatform, error)
	GetByChannelAndPlatform(ctx context.Context, channelID uuid.UUID, platform platform.Platform) (model.ChannelPlatform, error)
	GetByPlatformChannelID(ctx context.Context, platform platform.Platform, platformChannelID string) (model.ChannelPlatform, error)
	ListByChannelID(ctx context.Context, channelID uuid.UUID) ([]model.ChannelPlatform, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.ChannelPlatform, error)
	Patch(ctx context.Context, id uuid.UUID, input PatchInput) (model.ChannelPlatform, error)
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

type PatchInput struct {
	Enabled        *bool
	BotConfigPatch json.RawMessage
}

func (i PatchInput) Validate() error {
	if len(i.BotConfigPatch) == 0 {
		return nil
	}

	var patch map[string]json.RawMessage
	if err := json.Unmarshal(i.BotConfigPatch, &patch); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidBotConfigPatch, err)
	}
	if patch == nil {
		return ErrInvalidBotConfigPatch
	}

	return nil
}
