package overlays_kappagen

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/libs/repositories/overlays_kappagen/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.KappagenOverlay, error)
	Create(ctx context.Context, input CreateInput) (model.KappagenOverlay, error)
	Update(ctx context.Context, channelID string, input UpdateInput) (model.KappagenOverlay, error)
}

type CreateInput struct {
	ChannelID string
	Settings  model.KappagenOverlaySettings
}

type UpdateInput struct {
	Settings model.KappagenOverlaySettings
}

var ErrNotFound = fmt.Errorf("not found")
