package auth

import (
	"context"
	"fmt"

	model "github.com/twirapp/twir/libs/gomodels"
)

func (s *Auth) GetAuthenticatedUser(ctx context.Context) (*model.Users, error) {
	userByApyKey, err := s.GetAuthenticatedUserByApiKey(ctx)
	if err == nil {
		return userByApyKey, nil
	}

	user, ok := s.sessionManager.Get(ctx, "dbUser").(model.Users)
	if !ok {
		return nil, fmt.Errorf("not authenticated")
	}

	freshUser := model.Users{}
	if err := s.gorm.First(&freshUser, user.ID).Error; err != nil {
		return nil, fmt.Errorf("cannot get user from db: %w", err)
	}

	if freshUser.IsBanned {
		return nil, fmt.Errorf("forbidden")
	}

	return &freshUser, nil
}

func (s *Auth) GetSelectedDashboard(ctx context.Context) (string, error) {
	userByApyKey, err := s.GetAuthenticatedUserByApiKey(ctx)
	if err == nil {
		return userByApyKey.ID, nil
	}

	selectedDashboardId, ok := s.sessionManager.Get(ctx, "dashboardId").(string)
	if !ok {
		return "", fmt.Errorf("cannot get dashboardId from context")
	}

	return selectedDashboardId, nil
}

func (s *Auth) SetSessionSelectedDashboard(ctx context.Context, dashboardId string) error {
	s.sessionManager.Put(ctx, "dashboardId", dashboardId)
	s.sessionManager.Commit(ctx)

	return nil
}

func (s *Auth) SessionLogout(ctx context.Context) error {
	return s.sessionManager.Destroy(ctx)
}
