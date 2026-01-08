package channels_overlays

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/channels_overlays/model"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Overlay, error)
	GetManyByChannelID(ctx context.Context, channelID string) ([]model.Overlay, error)
	Create(ctx context.Context, input CreateInput) (model.Overlay, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.Overlay, error)
	UpdateLayer(ctx context.Context, layerId uuid.UUID, input LayerUpdateInput) (model.OverlayLayer, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateLayerInput struct {
	Type                    model.OverlayType
	Settings                model.OverlayLayerSettings
	PosX                    int
	PosY                    int
	Width                   int
	Height                  int
	Rotation                int
	PeriodicallyRefetchData bool
	Locked                  bool
	Visible                 bool
	Opacity                 float64
}

type CreateInput struct {
	ChannelID string
	Name      string
	Width     int
	Height    int
	InstaSave bool
	Layers    []CreateLayerInput
}

type UpdateLayerInputWithID struct {
	ID                      *uuid.UUID // nil for new layers, set for existing layers
	Type                    model.OverlayType
	Settings                model.OverlayLayerSettings
	PosX                    int
	PosY                    int
	Width                   int
	Height                  int
	Rotation                int
	PeriodicallyRefetchData bool
	Locked                  bool
	Visible                 bool
	Opacity                 float64
}

type UpdateInput struct {
	Name      string
	Width     int
	Height    int
	InstaSave bool
	Layers    []UpdateLayerInputWithID
}

type LayerUpdateInput struct {
	PosX     *int
	PosY     *int
	Width    *int
	Height   *int
	Rotation *int
	Visible  *bool
	Opacity  *float64
}
