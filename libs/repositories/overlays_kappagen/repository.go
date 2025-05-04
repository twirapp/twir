package overlays_kappagen

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/overlays_kappagen/model"
)

type Repository interface {
	// Get returns the kappagen overlay for the given channel ID
	Get(ctx context.Context, channelID string) (model.KappagenOverlay, error)
	
	// Create creates a new kappagen overlay for the given channel ID
	Create(ctx context.Context, input CreateInput) (model.KappagenOverlay, error)
	
	// Update updates the kappagen overlay for the given channel ID
	Update(ctx context.Context, channelID string, input UpdateInput) (model.KappagenOverlay, error)
}

type CreateInput struct {
	ChannelID      string
	EnableSpawn    bool
	ExcludedEmotes []string
	EnableRave     bool
	Animation      model.KappagenOverlayAnimationSettings
	Animations     []model.KappagenOverlayAnimationsSettings
	Emotes         model.KappagenOverlayEmotesSettings
	Size           model.KappagenOverlaySizeSettings
	Cube           model.KappagenOverlayCubeSettings
}

type UpdateInput struct {
	EnableSpawn    bool
	ExcludedEmotes []string
	EnableRave     bool
	Animation      model.KappagenOverlayAnimationSettings
	Animations     []model.KappagenOverlayAnimationsSettings
}

// ErrNotFound is returned when the kappagen overlay is not found
var ErrNotFound = &NotFoundError{}

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "kappagen overlay not found"
}
