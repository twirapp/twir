package short_links_custom_domains

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/short_links_custom_domains/model"
)

type Repository interface {
	GetByUserID(ctx context.Context, userID string) (model.CustomDomain, error)
	GetByDomain(ctx context.Context, domain string) (model.CustomDomain, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
	Create(ctx context.Context, input CreateInput) (model.CustomDomain, error)
	Update(ctx context.Context, id string, input UpdateInput) (model.CustomDomain, error)
	Delete(ctx context.Context, id string) error
	VerifyDomain(ctx context.Context, id string) error
}

type CreateInput struct {
	UserID            string
	Domain            string
	VerificationToken string
}

type UpdateInput struct {
	Domain            *string
	Verified          *bool
	VerificationToken *string
}
