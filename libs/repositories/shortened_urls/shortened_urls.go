package shortened_urls

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
)

type Repository interface {
	GetByShortID(ctx context.Context, id string) (model.ShortenedUrl, error)
	GetByUrl(ctx context.Context, url string) (model.ShortenedUrl, error)
	Create(ctx context.Context, input CreateInput) (model.ShortenedUrl, error)
	Update(ctx context.Context, id string, input UpdateInput) (model.ShortenedUrl, error)
}

type CreateInput struct {
	ShortID         string
	URL             string
	CreatedByUserID *string
}

type UpdateInput struct {
	Views *int
}
