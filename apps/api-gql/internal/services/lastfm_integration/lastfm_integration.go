package lastfmintegration

import (
	"context"
	"fmt"

	lfmapi "github.com/shkh/lastfm-go/lastfm"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	cfg "github.com/twirapp/twir/libs/config"
	channelsintegrationslastfm "github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	LastfmRepository channelsintegrationslastfm.Repository
	Config           cfg.Config
}

func New(opts Opts) *Service {
	return &Service{
		repo:   opts.LastfmRepository,
		config: opts.Config,
	}
}

type Service struct {
	repo   channelsintegrationslastfm.Repository
	config cfg.Config
}

func (s *Service) GetAuthLink(ctx context.Context) (string, error) {
	if s.config.LastFM.ApiKey == "" || s.config.LastFM.RedirectURL == "" {
		return "", fmt.Errorf("lastfm integration not configured")
	}

	link := fmt.Sprintf(
		"https://www.last.fm/api/auth/?api_key=%s&cb=%s",
		s.config.LastFM.ApiKey,
		s.config.LastFM.RedirectURL,
	)

	return link, nil
}

func (s *Service) GetData(
	ctx context.Context,
	channelID string,
) (*entity.LastfmIntegrationData, error) {
	integration, err := s.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lastfm data: %w", err)
	}

	if integration.IsNil() {
		return &entity.LastfmIntegrationData{
			Enabled: false,
		}, nil
	}

	result := &entity.LastfmIntegrationData{
		Enabled:  integration.Enabled,
		UserName: integration.UserName,
		Avatar:   integration.Avatar,
	}

	return result, nil
}

func (s *Service) PostCode(
	ctx context.Context,
	channelID string,
	code string,
) error {
	if s.config.LastFM.ApiKey == "" || s.config.LastFM.ClientSecret == "" {
		return fmt.Errorf("lastfm integration not configured")
	}

	// Create LastFM API client
	api := lfmapi.New(s.config.LastFM.ApiKey, s.config.LastFM.ClientSecret)

	// Exchange code for session
	err := api.LoginWithToken(code)
	if err != nil {
		return fmt.Errorf("failed to login with token: %w", err)
	}

	sessionKey := api.GetSessionKey()

	// Get user info
	userInfo, err := api.User.GetInfo(make(map[string]interface{}))
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	var avatar *string
	if len(userInfo.Images) > 0 {
		// Last.fm API returns images in ascending size order; select the largest (last) image as avatar.
		avatar = &userInfo.Images[len(userInfo.Images)-1].Url
	}

	userName := userInfo.Name

	// Check if integration already exists
	existing, err := s.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to check existing integration: %w", err)
	}

	if !existing.IsNil() {
		// Update existing integration
		err = s.repo.Update(
			ctx, existing.ID, channelsintegrationslastfm.UpdateInput{
				Enabled:    boolPtr(true),
				SessionKey: &sessionKey,
				UserName:   &userName,
				Avatar:     avatar,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to update lastfm integration: %w", err)
		}
	} else {
		// Create new integration
		_, err = s.repo.Create(
			ctx, channelsintegrationslastfm.CreateInput{
				ChannelID:  channelID,
				Enabled:    true,
				SessionKey: &sessionKey,
				UserName:   &userName,
				Avatar:     avatar,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create lastfm integration: %w", err)
		}
	}

	return nil
}

func (s *Service) Logout(
	ctx context.Context,
	channelID string,
) error {
	integration, err := s.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}

	if integration.IsNil() {
		return nil
	}

	return s.repo.Delete(ctx, integration.ID)
}

func boolPtr(b bool) *bool {
	return &b
}
