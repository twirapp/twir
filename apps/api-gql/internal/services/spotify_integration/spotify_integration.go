package spotify_integration

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/imroc/req/v3"
	deprecatedgormmodel "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/integrations/spotify"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	"github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/model"
	"github.com/twirapp/twir/libs/repositories/integrations"
	integrationmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

type spotifyTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type Service struct {
	spotifyRepository channelsintegrationsspotify.Repository
	integrationsRepo  integrations.Repository
}

func New(
	spotifyRepository channelsintegrationsspotify.Repository,
	integrationsRepo integrations.Repository,
) *Service {
	return &Service{
		spotifyRepository: spotifyRepository,
		integrationsRepo:  integrationsRepo,
	}
}

func (s *Service) GetSpotifyData(
	ctx context.Context,
	channelID string,
) (*entity.SpotifyIntegrationData, error) {
	integration, err := s.spotifyRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spotify data: %w", err)
	}

	if integration.ID == uuid.Nil {
		return nil, nil
	}

	return &entity.SpotifyIntegrationData{
		UserName: integration.Username,
		Avatar:   integration.AvatarURI,
	}, nil
}

func (s *Service) PostCode(
	ctx context.Context,
	channelID string,
	code string,
) error {
	integration, err := s.integrationsRepo.GetByService(ctx, integrationmodel.ServiceSpotify)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}

	if integration.ClientID == nil || integration.ClientSecret == nil || integration.RedirectURL == nil {
		return fmt.Errorf("spotify not enabled on our side, please be patient")
	}

	var data spotifyTokensResponse
	resp, err := req.R().
		SetContext(ctx).
		SetFormData(
			map[string]string{
				"grant_type":   "authorization_code",
				"redirect_uri": *integration.RedirectURL,
				"code":         code,
			},
		).
		SetBasicAuth(
			*integration.ClientID,
			*integration.ClientSecret,
		).
		SetSuccessResult(&data).
		SetContentType("application/x-www-form-urlencoded").
		Post("https://accounts.spotify.com/api/token")
	if err != nil {
		return fmt.Errorf("failed to get spotify tokens: %w", err)
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("failed to get spotify tokens: %s", resp.String())
	}

	createInput := channelsintegrationsspotify.CreateInput{
		ChannelID:    channelID,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		Scopes:       []string{"user-read-currently-playing", "user-read-playback-state"},
	}

	userSpotify := spotify.New(
		deprecatedgormmodel.Integrations{
			ID:      "",
			Service: "",
			AccessToken: sql.NullString{
				String: data.AccessToken,
				Valid:  true,
			},
			RefreshToken: sql.NullString{
				String: data.RefreshToken,
				Valid:  true,
			},
			ClientID: sql.NullString{
				String: *integration.ClientID,
				Valid:  true,
			},
			ClientSecret: sql.NullString{
				String: *integration.ClientSecret,
				Valid:  true,
			},
			APIKey: sql.NullString{},
			RedirectURL: sql.NullString{
				String: *integration.RedirectURL,
				Valid:  true,
			},
		},
		model.ChannelIntegrationSpotify{
			AccessToken:  data.AccessToken,
			RefreshToken: data.RefreshToken,
		},
		s.spotifyRepository,
	)
	profile, err := userSpotify.GetProfile(ctx)
	if err != nil {
		return fmt.Errorf("failed to get spotify profile: %w", err)
	}

	createInput.AvatarURI = profile.Images[0].URL
	createInput.Username = profile.DisplayName

	if _, err := s.spotifyRepository.Create(ctx, createInput); err != nil {
		return fmt.Errorf("failed to create spotify integration: %w", err)
	}

	return nil
}

func (s *Service) Logout(
	ctx context.Context,
	channelID string,
) error {
	integration, err := s.spotifyRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return err
	}

	return s.spotifyRepository.Delete(ctx, integration.ID)
}

func (s *Service) GetAuthLink(
	ctx context.Context,
) (string, error) {
	integration, err := s.integrationsRepo.GetByService(ctx, integrationmodel.ServiceSpotify)
	if err != nil {
		return "", err
	}

	if integration.ClientID == nil || integration.ClientSecret == nil || integration.RedirectURL == nil {
		return "", fmt.Errorf("spotify not enabled on our side, please be patient")
	}

	link, _ := url.Parse("https://accounts.spotify.com/authorize")

	q := link.Query()
	q.Add("response_type", "code")
	q.Add("client_id", *integration.ClientID)
	q.Add("scope", "user-read-currently-playing user-read-playback-state user-read-recently-played")
	q.Add("redirect_uri", *integration.RedirectURL)
	link.RawQuery = q.Encode()

	return link.String(), nil
}
