package tokens

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/tokens/model"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Token, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Token, error)
	GetByBotID(ctx context.Context, botID string) (*model.Token, error)
	CreateUserToken(ctx context.Context, input CreateInput) (*model.Token, error)
	UpdateTokenByID(ctx context.Context, id uuid.UUID, input UpdateTokenInput) (
		*model.Token,
		error,
	)
}

type CreateInput struct {
	UserID              uuid.UUID
	AccessToken         string
	RefreshToken        string
	ExpiresIn           int
	ObtainmentTimestamp time.Time
	Scopes              []string
}

type UpdateTokenInput struct {
	AccessToken         *string
	RefreshToken        *string
	ExpiresIn           *int
	ObtainmentTimestamp *time.Time
	Scopes              []string
}
