package short_links_banned_ua_preset_patterns

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("pattern not found")
var ErrAlreadyExists = errors.New("pattern already exists in this preset")

type Repository interface {
	GetByPresetID(ctx context.Context, presetID string) ([]Pattern, error)
	Create(ctx context.Context, input CreateInput) (Pattern, error)
	Delete(ctx context.Context, id string, presetID string) error
}

type Pattern struct {
	ID          string
	PresetID    string
	Pattern     string
	Description *string
	CreatedAt   time.Time

	isNil bool
}

func (p Pattern) IsNil() bool {
	return p.isNil
}

var Nil = Pattern{isNil: true}

type CreateInput struct {
	PresetID    string
	Pattern     string
	Description *string
}
