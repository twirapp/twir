package auth

import (
	"context"
	"fmt"

	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
)

func (s *Auth) GetAuthenticatedUserByApiKey(ctx context.Context) (*model.Users, error) {
	var apiKey string

	wsApiKey, _ := s.getWsAuthenticatedApiKey(ctx)
	if wsApiKey != "" {
		apiKey = wsApiKey
	} else {
		ginCtx, err := gincontext.GetGinContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get gin context: %w", err)
		}

		apiKey = ginCtx.GetHeader("api-key")
	}

	if apiKey == "" {
		return nil, fmt.Errorf("api key is required")
	}

	// Try channel API key first
	channel := model.Channels{}
	if err := s.gorm.Where(`"api_key" = ?`, apiKey).First(&channel).Error; err == nil {
		// Found channel by API key, resolve owner user
		var userID string
		if channel.TwitchUserID != nil {
			userID = *channel.TwitchUserID
		} else if channel.KickUserID != nil {
			userID = *channel.KickUserID
		}
		if userID != "" {
			user := model.Users{}
			if err := s.gorm.Where("id = ?", userID).First(&user).Error; err == nil {
				return &user, nil
			}
		}
	}

	// Fallback to user API key
	user := model.Users{}
	if err := s.gorm.Where(`"apiKey" = ?`, apiKey).First(&user).Error; err != nil {
		return nil, fmt.Errorf("cannot get user from db: %w", err)
	}

	return &user, nil
}
