package short_links_link_banned_user_agents

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("banned user agent not found")
var ErrAlreadyExists = errors.New("banned user agent pattern already exists")

type Repository interface {
	GetByLinkID(ctx context.Context, linkID string) ([]BannedUserAgent, error)
	Create(ctx context.Context, input CreateInput) (BannedUserAgent, error)
	Delete(ctx context.Context, id string, linkID string) error
}

type BannedUserAgent struct {
	ID          string
	LinkID      string
	Pattern     string
	Description *string
	CreatedAt   time.Time

	isNil bool
}

func (b BannedUserAgent) IsNil() bool {
	return b.isNil
}

var Nil = BannedUserAgent{isNil: true}

type CreateInput struct {
	LinkID      string
	Pattern     string
	Description *string
}
