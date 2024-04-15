package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/satont/twir/apps/api/internal/helpers"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/users"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Users struct {
	*impl_deps.Deps
}

func (c *Users) UsersRegenerateApiKey(
	ctx context.Context,
	req *users.RegenerateApiKeyRequest,
) (*emptypb.Empty, error) {
	user, err := helpers.GetUserModelFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from context: %w", err)
	}

	user.ApiKey = uuid.New().String()

	if err := c.Db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (c *Users) UsersUpdate(ctx context.Context, req *users.UpdateUserRequest) (
	*emptypb.Empty,
	error,
) {
	requestUser, err := helpers.GetUserModelFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from context: %w", err)
	}

	user := model.Users{}
	if err := c.Db.WithContext(ctx).Where("id = ?", requestUser.ID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}

	query := c.Db.WithContext(ctx).Model(&user)
	// everything bellow working like a PATCH http request
	if req.HideOnLandingPage != nil {
		user.HideOnLandingPage = *req.HideOnLandingPage
	}

	if err := query.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	c.SessionManager.Put(ctx, "dbUser", user)

	return &emptypb.Empty{}, nil
}
