package auth

import (
	"context"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	model "github.com/twirapp/twir/libs/gomodels"
)

const (
	latestShortenedUrlsIdsKey = "latestShortenedUrlsIds"
	dbUserKey                 = "dbUser"
	dashboardIdKey            = "dashboardId"
	twitchUserKey             = "twitchUser"
)

func (s *Auth) GetLatestShortenerUrlsIds(ctx context.Context) ([]string, error) {
	ids, ok := s.sessionManager.Get(ctx, latestShortenedUrlsIdsKey).([]string)
	if !ok {
		return nil, fmt.Errorf("not authenticated")
	}

	return ids, nil
}

func (s *Auth) AddLatestShortenerUrlsId(ctx context.Context, id string) error {
	latest, _ := s.GetLatestShortenerUrlsIds(ctx)
	latest = append([]string{id}, latest...)
	if len(latest) > 5 {
		latest = latest[:5]
	}

	s.sessionManager.Put(ctx, latestShortenedUrlsIdsKey, latest)
	if _, _, err := s.sessionManager.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit session: %w", err)
	}

	return nil
}

func (s *Auth) GetAuthenticatedUserModel(ctx context.Context) (*model.Users, error) {
	userByApyKey, err := s.GetAuthenticatedUserByApiKey(ctx)
	if err == nil {
		return userByApyKey, nil
	}

	user, ok := s.sessionManager.Get(ctx, dbUserKey).(model.Users)
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

func (s *Auth) SetSessionAuthenticatedUser(ctx context.Context, user model.Users) error {
	s.sessionManager.Put(ctx, dbUserKey, user)
	if _, _, err := s.sessionManager.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit session: %w", err)
	}

	return nil
}

func (s *Auth) SetSessionTwitchUser(ctx context.Context, user helix.User) error {
	s.sessionManager.Put(ctx, twitchUserKey, user)
	if _, _, err := s.sessionManager.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit session: %w", err)
	}

	return nil
}

func (s *Auth) GetSelectedDashboard(ctx context.Context) (string, error) {
	userByApyKey, err := s.GetAuthenticatedUserByApiKey(ctx)
	if err == nil {
		return userByApyKey.ID, nil
	}

	selectedDashboardId, ok := s.sessionManager.Get(ctx, dashboardIdKey).(string)
	if !ok {
		return "", fmt.Errorf("cannot get dashboardId from context")
	}

	return selectedDashboardId, nil
}

func (s *Auth) SetSessionSelectedDashboard(ctx context.Context, dashboardId string) error {
	s.sessionManager.Put(ctx, dashboardIdKey, dashboardId)
	if _, _, err := s.sessionManager.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit session: %w", err)
	}

	return nil
}

func (s *Auth) SessionLogout(ctx context.Context) error {
	s.sessionManager.Remove(ctx, dbUserKey)
	s.sessionManager.Remove(ctx, dashboardIdKey)
	s.sessionManager.Remove(ctx, "twitchUser")
	if _, _, err := s.sessionManager.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit session: %w", err)
	}

	fmt.Println(s.sessionManager.Keys(ctx))

	return nil
}
