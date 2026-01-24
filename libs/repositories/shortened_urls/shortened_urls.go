package shortened_urls

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
)

type Repository interface {
	GetByShortID(ctx context.Context, domain *string, id string) (model.ShortenedUrl, error)
	GetManyByShortIDs(ctx context.Context, domain *string, ids []string) ([]model.ShortenedUrl, error)
	GetByUrl(ctx context.Context, domain *string, url string) (model.ShortenedUrl, error)
	Create(ctx context.Context, input CreateInput) (model.ShortenedUrl, error)
	Update(ctx context.Context, domain *string, id string, input UpdateInput) (model.ShortenedUrl, error)
	GetList(ctx context.Context, input GetListInput) (GetListOutput, error)
	Delete(ctx context.Context, domain *string, id string) error
	Count(ctx context.Context, input CountInput) (int64, error)
}

type CreateInput struct {
	ShortID         string
	URL             string
	CreatedByUserID *string
	UserIp          *string
	UserAgent       *string
	Domain          *string
}

type UpdateInput struct {
	Views   *int
	ShortID *string
	URL     *string
	Domain  *string
}

type GetListInput struct {
	Page    int
	PerPage int
	UserID  *string
	SortBy  string // "views" or "created_at"
}

type GetListOutput struct {
	Items []model.ShortenedUrl
	Total int
}

type CountInput struct {
	UserID string
}
