package pastebins

import (
	"context"
	"time"

	"github.com/twirapp/twir/libs/repositories/pastebins/model"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) (model.Pastebin, error)
	GetByID(ctx context.Context, id string) (model.Pastebin, error)
	Delete(ctx context.Context, id string) error
	GetManyByOwner(ctx context.Context, input GetManyInput) (GetManyOutput, error)
}

type CreateInput struct {
	ID          string
	Content     string
	ExpireAt    *time.Time
	OwnerUserID *string
}

type GetManyInput struct {
	Page        int
	PerPage     int
	OwnerUserID string
}

type GetManyOutput struct {
	Items []model.Pastebin
	Total int
}
