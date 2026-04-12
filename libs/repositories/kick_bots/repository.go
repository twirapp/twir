package kick_bots

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	entity "github.com/twirapp/twir/libs/entities/kick_bot"
)

type Repository interface {
	GetDefault(ctx context.Context) (entity.KickBot, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.KickBot, error)
	Create(ctx context.Context, input CreateInput) (entity.KickBot, error)
	UpdateToken(ctx context.Context, id uuid.UUID, input UpdateTokenInput) (entity.KickBot, error)
}

type CreateInput struct {
	Type                string
	AccessToken         string
	RefreshToken        string
	Scopes              []string
	ExpiresIn           int
	ObtainmentTimestamp time.Time
	KickUserID          string
	KickUserLogin       string
}

type UpdateTokenInput struct {
	AccessToken         string
	RefreshToken        string
	Scopes              []string
	ExpiresIn           int
	ObtainmentTimestamp time.Time
}

var ErrNotFound = fmt.Errorf("kick bot not found")
