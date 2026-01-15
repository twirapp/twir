package short_links_views

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
}

type CreateInput struct {
	ShortLinkID string
	UserID      *string
	IP          *string
	UserAgent   *string
	CreatedAt   time.Time
}
