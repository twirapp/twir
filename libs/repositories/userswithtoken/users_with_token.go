package userswithtoken

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/userswithtoken/model"
)

type Repository interface {
	GetByID(ctx context.Context, userID string) (model.UserWithToken, error)
}
