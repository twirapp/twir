package vk_video_bots

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	entity "github.com/twirapp/twir/libs/entities/vk_video_bot"
)

var ErrNotFound = errors.New("VK Video bot not found")

type Repository interface {
	Get(ctx context.Context) (entity.VKVideoBot, error)
	Lock(ctx context.Context) error
	Upsert(ctx context.Context, input UpsertInput) (entity.VKVideoBot, error)
	Update(ctx context.Context, input UpdateInput) (entity.VKVideoBot, error)
}

type UpsertInput struct {
	EncryptedAccessToken  string
	EncryptedRefreshToken string
	Scopes                []string
	ExpiresIn             int
	ObtainmentTimestamp   time.Time
	VKUserID              uuid.UUID
}

type UpdateInput struct {
	EncryptedAccessToken  string
	EncryptedRefreshToken string
	Scopes                []string
	ExpiresIn             int
	ObtainmentTimestamp   time.Time
	VKUserID              uuid.UUID
}
