package short_links_link_presets

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("link preset not found")
var ErrAlreadyExists = errors.New("preset already applied to this link")

type Repository interface {
	GetByLinkID(ctx context.Context, linkID string) ([]LinkPreset, error)
	GetByPresetID(ctx context.Context, presetID string) ([]LinkPreset, error)
	GetLinksByPresetID(ctx context.Context, presetID string) ([]string, error)
	Create(ctx context.Context, input CreateInput) (LinkPreset, error)
	Delete(ctx context.Context, id string) error
	DeleteByLinkAndPreset(ctx context.Context, linkID string, presetID string) error
}

type LinkPreset struct {
	ID        string
	LinkID    string
	PresetID  string
	CreatedAt time.Time

	isNil bool
}

func (l LinkPreset) IsNil() bool {
	return l.isNil
}

var Nil = LinkPreset{isNil: true}

type CreateInput struct {
	LinkID   string
	PresetID string
}
