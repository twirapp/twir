package user_platform_accounts

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	entity "github.com/twirapp/twir/libs/entities/user_platform_account"
)

type Repository interface {
	GetByUserIDAndPlatform(ctx context.Context, userID uuid.UUID, platform platform.Platform) (entity.UserPlatformAccount, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]entity.UserPlatformAccount, error)
	GetByPlatformUserID(ctx context.Context, plat platform.Platform, platformUserID string) (entity.UserPlatformAccount, error)
	Upsert(ctx context.Context, input UpsertInput) (entity.UserPlatformAccount, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type UpsertInput struct {
	UserID              uuid.UUID
	Platform            platform.Platform
	PlatformUserID      string
	PlatformLogin       string
	PlatformDisplayName string
	PlatformAvatar      string
	AccessToken         string
	RefreshToken        string
	Scopes              []string
	ExpiresIn           int
	ObtainmentTimestamp time.Time
}

var ErrNotFound = fmt.Errorf("user platform account not found")
