package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	model "github.com/twirapp/twir/libs/gomodels"
)

const (
	latestShortenedUrlsIdsKey = "latestShortenedUrlsIds"
	dbUserKey                 = "dbUser"
	dashboardIdKey            = "dashboardId"
	twitchUserKey             = "twitchUser"
	internalUserIdKey         = "internalUserId"
	currentPlatformKey        = "currentPlatform"
	selectedDashboardIdKey    = "selectedDashboardId"
	kickUserKey               = "kickUser"
)

type KickSessionUser struct {
	ID     string
	Login  string
	Avatar string
}

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

	userID, err := s.GetInternalUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("not authenticated")
	}

	dbUser, err := s.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cannot get user from db: %w", err)
	}

	freshUser := &model.Users{
		ID:                dbUser.ID.String(),
		TokenID:           dbUser.TokenID.NullString,
		IsBotAdmin:        dbUser.IsBotAdmin,
		ApiKey:            dbUser.ApiKey,
		IsBanned:          dbUser.IsBanned,
		HideOnLandingPage: dbUser.HideOnLandingPage,
		CreatedAt:         dbUser.CreatedAt,
	}

	if freshUser.IsBanned {
		return nil, fmt.Errorf("forbidden")
	}

	return freshUser, nil
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
	s.sessionManager.Remove(ctx, internalUserIdKey)
	s.sessionManager.Remove(ctx, currentPlatformKey)
	s.sessionManager.Remove(ctx, selectedDashboardIdKey)
	s.sessionManager.Remove(ctx, kickUserKey)
	if _, _, err := s.sessionManager.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit session: %w", err)
	}

	fmt.Println(s.sessionManager.Keys(ctx))

	return nil
}

func (s *Auth) SetSessionInternalUserID(ctx context.Context, id uuid.UUID) error {
	s.sessionManager.Put(ctx, internalUserIdKey, id)
	if _, _, err := s.sessionManager.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit session: %w", err)
	}

	return nil
}

func (s *Auth) GetInternalUserID(ctx context.Context) (uuid.UUID, error) {
	id, ok := s.sessionManager.Get(ctx, internalUserIdKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("internalUserId not found in session")
	}

	return id, nil
}

func (s *Auth) SetSessionCurrentPlatform(ctx context.Context, platform string) error {
	s.sessionManager.Put(ctx, currentPlatformKey, platform)
	if _, _, err := s.sessionManager.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit session: %w", err)
	}

	return nil
}

func (s *Auth) GetCurrentPlatform(ctx context.Context) (string, error) {
	platform, ok := s.sessionManager.Get(ctx, currentPlatformKey).(string)
	if !ok {
		return "", fmt.Errorf("currentPlatform not found in session")
	}

	return platform, nil
}

func (s *Auth) SetSelectedDashboardUUID(ctx context.Context, channelID uuid.UUID) error {
	s.sessionManager.Put(ctx, selectedDashboardIdKey, channelID)
	if _, _, err := s.sessionManager.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit session: %w", err)
	}

	return nil
}

func (s *Auth) GetSelectedDashboardUUID(ctx context.Context) (uuid.UUID, error) {
	id, ok := s.sessionManager.Get(ctx, selectedDashboardIdKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("selectedDashboardId not found in session")
	}

	return id, nil
}

func (s *Auth) SetSessionKickUser(ctx context.Context, user KickSessionUser) error {
	s.sessionManager.Put(ctx, kickUserKey, user)
	if _, _, err := s.sessionManager.Commit(ctx); err != nil {
		return fmt.Errorf("cannot commit session: %w", err)
	}

	return nil
}

func (s *Auth) GetSessionKickUser(ctx context.Context) (KickSessionUser, error) {
	user, ok := s.sessionManager.Get(ctx, kickUserKey).(KickSessionUser)
	if !ok {
		return KickSessionUser{}, fmt.Errorf("kickUser not found in session")
	}

	return user, nil
}
