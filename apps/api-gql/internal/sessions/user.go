package sessions

import (
	"context"
	"fmt"

	model "github.com/satont/twir/libs/gomodels"
)

func (s *Sessions) GetAuthenticatedUser(ctx context.Context) (*model.Users, error) {
	user, ok := s.sessionManager.Get(ctx, "dbUser").(model.Users)
	if !ok {
		return nil, fmt.Errorf("not authenticated")
	}

	freshUser := model.Users{}
	if err := s.gorm.First(&freshUser, user.ID).Error; err != nil {
		return nil, fmt.Errorf("cannot get user from db: %w", err)
	}

	return &freshUser, nil
}

func (s *Sessions) GetSelectedDashboard(ctx context.Context) (string, error) {
	selectedDashboardId, ok := s.sessionManager.Get(ctx, "dashboardId").(string)
	if !ok {
		return "", fmt.Errorf("cannot get dashboardId from context")
	}

	return selectedDashboardId, nil
}

func (s *Sessions) SetSelectedDashboard(ctx context.Context, dashboardId string) error {
	s.sessionManager.Put(ctx, "dashboardId", dashboardId)
	s.sessionManager.Commit(ctx)

	return nil
}

func (s *Sessions) Logout(ctx context.Context) error {
	return s.sessionManager.Destroy(ctx)
}
