package short_links_global_banned_user_agents

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("banned user agent not found")
var ErrAlreadyExists = errors.New("banned user agent pattern already exists")

type Repository interface {
	GetByUserID(ctx context.Context, userID string) ([]BannedUserAgent, error)
	Create(ctx context.Context, input CreateInput) (BannedUserAgent, error)
	Delete(ctx context.Context, id string, userID string) error
}

type BannedUserAgent struct {
	ID          string
	UserID      string
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
	UserID      string
	Pattern     string
	Description *string
}
