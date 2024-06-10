package auth

import (
	"context"
	"fmt"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gincontext"
)

func (s *Auth) GetAuthenticatedUserByApiKey(ctx context.Context) (*model.Users, error) {
	ginCtx, err := gincontext.GetGinContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gin context: %w", err)
	}

	apiKey := ginCtx.GetHeader("x-api-key")
	if apiKey == "" {
		return nil, fmt.Errorf("api key is required")
	}

	user := model.Users{}
	if err := s.gorm.Where(`"apiKey" = ?`, apiKey).First(&user).Error; err != nil {
		return nil, fmt.Errorf("cannot get user from db: %w", err)
	}

	return &user, nil
}
