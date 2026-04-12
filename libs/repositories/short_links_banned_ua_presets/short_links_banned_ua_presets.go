package short_links_banned_ua_presets

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("preset not found")
var ErrAlreadyExists = errors.New("preset with this name already exists")

type Repository interface {
	GetByUserID(ctx context.Context, userID string) ([]Preset, error)
	GetByID(ctx context.Context, id string) (Preset, error)
	Create(ctx context.Context, input CreateInput) (Preset, error)
	Update(ctx context.Context, id string, input UpdateInput) (Preset, error)
	Delete(ctx context.Context, id string, userID string) error
}

type Preset struct {
	ID          string
	UserID      string
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	isNil bool
}

func (p Preset) IsNil() bool {
	return p.isNil
}

var Nil = Preset{isNil: true}

type CreateInput struct {
	UserID      string
	Name        string
	Description *string
}

type UpdateInput struct {
	Name        *string
	Description *string
}
