package overlays_tts

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/libs/repositories/overlays_tts/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.TTSOverlay, error)
	Create(ctx context.Context, input CreateInput) (model.TTSOverlay, error)
	Update(ctx context.Context, channelID string, input UpdateInput) (model.TTSOverlay, error)
}

type CreateInput struct {
	ChannelID string
	Settings  model.TTSOverlaySettings
}

type UpdateInput struct {
	Settings model.TTSOverlaySettings
}

var ErrNotFound = fmt.Errorf("not found")

