package spotify_integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	config "github.com/twirapp/twir/libs/config"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
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
	config            config.Config
}

func New(
	spotifyRepository channelsintegrationsspotify.Repository,
	integrationsRepo integrations.Repository,
	cfg config.Config,
) *Service {
	return &Service{
		spotifyRepository: spotifyRepository,
		integrationsRepo:  integrationsRepo,
		config:            cfg,
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

func (s *Service) getCallbackUrl(ctx context.Context) (string, error) {
	baseUrl, _ := gincontext.GetBaseUrlFromContext(ctx, s.config.SiteBaseUrl)
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("invalid site base URL: %w", err)
	}

	return u.JoinPath("dashboard", "integrations", "spotify").String(), nil
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

	if integration.ClientID == nil || integration.ClientSecret == nil {
		return fmt.Errorf("spotify not enabled on our side, please be patient")
	}

	redirectUrl, err := s.getCallbackUrl(ctx)
	if err != nil {
		return fmt.Errorf("failed to get redirect URL: %w", err)
	}

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("redirect_uri", redirectUrl)
	formData.Set("code", code)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://accounts.spotify.com/api/token",
		bytes.NewBufferString(formData.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString(
		[]byte(*integration.ClientID + ":" + *integration.ClientSecret),
	)
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get spotify tokens: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to get spotify tokens: %s", string(body))
	}

	var data spotifyTokensResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
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
			APIKey:      sql.NullString{},
			RedirectURL: sql.NullString{},
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

	if len(profile.Images) > 0 {
		createInput.AvatarURI = profile.Images[0].URL
	}

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

	if integration.ClientID == nil || integration.ClientSecret == nil {
		return "", fmt.Errorf("spotify not enabled on our side, please be patient")
	}

	redirectUrl, err := s.getCallbackUrl(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get redirect URL: %w", err)
	}

	link, _ := url.Parse("https://accounts.spotify.com/authorize")

	q := link.Query()
	q.Add("response_type", "code")
	q.Add("client_id", *integration.ClientID)
	q.Add("scope", "user-read-currently-playing user-read-playback-state user-read-recently-played")
	q.Add("redirect_uri", redirectUrl)
	link.RawQuery = q.Encode()

	return link.String(), nil
}
