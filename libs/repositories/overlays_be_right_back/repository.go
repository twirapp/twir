package overlays_be_right_back

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/libs/repositories/overlays_be_right_back/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.BeRightBackOverlay, error)
	Create(ctx context.Context, input CreateInput) (model.BeRightBackOverlay, error)
	Update(ctx context.Context, channelID string, input UpdateInput) (model.BeRightBackOverlay, error)
}

type CreateInput struct {
	ChannelID string
	Settings  model.BeRightBackOverlaySettings
}

type UpdateInput struct {
	Settings model.BeRightBackOverlaySettings
}

var ErrNotFound = fmt.Errorf("not found")
