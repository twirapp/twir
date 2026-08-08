package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	"github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func (s *Auth) getApiKey(ctx context.Context) (string, error) {
	wsApiKey, _ := s.getWsAuthenticatedApiKey(ctx)
	if wsApiKey != "" {
		return wsApiKey, nil
	}

	ginCtx, err := gincontext.GetGinContext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gin context: %w", err)
	}

	apiKey := ginCtx.GetHeader("api-key")
	if apiKey == "" {
		return "", fmt.Errorf("api key is required")
	}

	return apiKey, nil
}

func (s *Auth) GetChannelFromApiKey(ctx context.Context) (channelentity.Channel, error) {
	apiKey, err := s.getApiKey(ctx)
	if err != nil {
		return channelentity.Nil, err
	}

	channel, err := s.channelService.GetChannelByApiKey(ctx, apiKey)
	if err != nil {
		return channelentity.Nil, fmt.Errorf("cannot get channel by api key: %w", err)
	}
	if channel.IsNil() {
		return channelentity.Nil, fmt.Errorf("cannot get channel by api key: channel not found")
	}

	return channel, nil
}

func (s *Auth) GetAuthenticatedUserByApiKey(ctx context.Context) (*model.Users, error) {
	var apiKey string
	var ginCtx *gin.Context

	wsApiKey, _ := s.getWsAuthenticatedApiKey(ctx)
	if wsApiKey != "" {
		apiKey = wsApiKey
	} else {
		var err error
		ginCtx, err = gincontext.GetGinContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get gin context: %w", err)
		}

		apiKey = ginCtx.GetHeader("api-key")
	}

	if apiKey == "" {
		return nil, fmt.Errorf("api key is required")
	}

	if user, found := s.resolveChannelAPIKeyOwner(ctx, apiKey); found {
		return user, nil
	}

	user, err := s.usersRepo.GetByApiKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("cannot get user from db: %w", err)
	}
	if user.IsNil() {
		return nil, fmt.Errorf("cannot get user from db: %w", usersmodel.ErrNotFound)
	}

	slog.WarnContext(ctx, "user API key is deprecated, use channel API key")
	if ginCtx != nil {
		ginCtx.Header("Deprecation", "true")
	}

	return mapRepositoryUser(user), nil
}

func (s *Auth) resolveChannelAPIKeyOwner(ctx context.Context, apiKey string) (*model.Users, bool) {
	channel, err := s.channelService.GetChannelByApiKey(ctx, apiKey)
	if err != nil || channel.IsNil() {
		return nil, false
	}

	for _, p := range platform.All() {
		for _, binding := range channel.Bindings {
			if binding.Platform != p {
				continue
			}

			return s.getBindingOwner(ctx, binding.UserID)
		}
	}

	for _, binding := range channel.Bindings {
		if binding.UserID == uuid.Nil {
			continue
		}

		return s.getBindingOwner(ctx, binding.UserID)
	}

	return nil, false
}

func (s *Auth) getBindingOwner(ctx context.Context, userID uuid.UUID) (*model.Users, bool) {
	if userID == uuid.Nil {
		return nil, false
	}

	user, err := s.usersRepo.GetByID(ctx, userID)
	if err != nil || user.IsNil() {
		return nil, false
	}

	return mapRepositoryUser(user), true
}

func mapRepositoryUser(user usersmodel.User) *model.Users {
	return &model.Users{
		ID:                user.ID.String(),
		PlatformID:        user.PlatformID,
		TokenID:           user.TokenID.NullString,
		IsBotAdmin:        user.IsBotAdmin,
		ApiKey:            user.ApiKey,
		IsBanned:          user.IsBanned,
		CreatedAt:         user.CreatedAt,
		HideOnLandingPage: user.HideOnLandingPage,
	}
}
