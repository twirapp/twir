package badges

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/badges/model"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Badge, error)
	GetMany(ctx context.Context, input GetManyInput) ([]model.Badge, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Create(ctx context.Context, input CreateInput) (model.Badge, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.Badge, error)
}

type GetManyInput struct {
	Enabled *bool
}

type CreateInput struct {
	// ID needed for passing into create for correct generate file name
	ID       uuid.UUID
	Name     string
	Enabled  bool
	FFZSlot  int
	FileName string
}

type UpdateInput struct {
	Name     *string
	Enabled  *bool
	FFZSlot  *int
	FileName *string
}
