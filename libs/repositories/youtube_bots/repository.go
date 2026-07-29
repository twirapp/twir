package youtube_bots

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	entity "github.com/twirapp/twir/libs/entities/youtube_bot"
)

var ErrNotFound = errors.New("YouTube bot not found")

type Repository interface {
	Get(ctx context.Context) (entity.YouTubeBot, error)
	Lock(ctx context.Context) error
	Upsert(ctx context.Context, input UpsertInput) (entity.YouTubeBot, error)
	Update(ctx context.Context, input UpdateInput) (entity.YouTubeBot, error)
}

type UpsertInput struct {
	EncryptedAccessToken  string
	EncryptedRefreshToken string
	Scopes                []string
	ExpiresIn             int
	ObtainmentTimestamp   time.Time
	YouTubeUserID         uuid.UUID
}

type UpdateInput struct {
	EncryptedAccessToken  string
	EncryptedRefreshToken string
	Scopes                []string
	ExpiresIn             int
	ObtainmentTimestamp   time.Time
	YouTubeUserID         uuid.UUID
}
