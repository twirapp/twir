package streamlabs_integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/integrations"
	config "github.com/twirapp/twir/libs/config"
	streamlabsintegration "github.com/twirapp/twir/libs/entities/streamlabs_integration"
	"github.com/twirapp/twir/libs/repositories/streamlabs_integration"
	"github.com/twirapp/twir/libs/repositories/streamlabs_integration/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	StreamlabsRepository streamlabs_integration.Repository
	TwirBus              *buscore.Bus
	Config               config.Config
}

func New(opts Opts) *Service {
	return &Service{
		streamlabsRepository: opts.StreamlabsRepository,
		twirBus:              opts.TwirBus,
		config:               opts.Config,
	}
}

type Service struct {
	streamlabsRepository streamlabs_integration.Repository
	twirBus              *buscore.Bus
	config               config.Config
}

type AuthLinkResponse struct {
	Link string `json:"link"`
}

// mapModelToEntity converts repository model to service entity
func (s *Service) mapModelToEntity(m model.StreamlabsIntegration) streamlabsintegration.Entity {
	return streamlabsintegration.Entity{
		ID:           m.ID,
		Enabled:      m.Enabled,
		ChannelID:    m.ChannelID,
		AccessToken:  m.AccessToken,
		RefreshToken: m.RefreshToken,
		UserName:     m.UserName,
		Avatar:       m.Avatar,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (s *Service) GetIntegrationData(ctx context.Context, channelID string) (
	streamlabsintegration.Entity,
	error,
) {
	integration, err := s.streamlabsRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, streamlabs_integration.ErrNotFound) {
			// Return default data if integration doesn't exist
			return streamlabsintegration.Entity{
				ChannelID: channelID,
				Enabled:   false,
			}, nil
		}
		return streamlabsintegration.Entity{}, fmt.Errorf(
			"failed to get streamlabs integration: %w",
			err,
		)
	}

	return s.mapModelToEntity(integration), nil
}

func (s *Service) getCallbackUrl(ctx context.Context) (string, error) {
	baseUrl, _ := gincontext.GetBaseUrlFromContext(ctx, s.config.SiteBaseUrl)
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("invalid site base URL: %w", err)
	}

	return u.JoinPath("dashboard", "integrations", "streamlabs").String(), nil
}

func (s *Service) GetAuthLink(ctx context.Context) (*AuthLinkResponse, error) {
	if s.config.StreamlabsClientId == "" || s.config.StreamlabsClientSecret == "" {
		return nil, errors.New("streamlabs integration not properly configured")
	}

	redirectUrl, err := s.getCallbackUrl(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get redirect URL: %w", err)
	}

	// Build OAuth authorization URL
	authURL := "https://www.streamlabs.com/api/v2.0/authorize"
	params := url.Values{}
	params.Add("client_id", s.config.StreamlabsClientId)
	params.Add("redirect_uri", redirectUrl)
	params.Add("response_type", "code")
	params.Add("scope", "socket.token donations.read")

	fullURL := fmt.Sprintf("%s?%s", authURL, params.Encode())

	return &AuthLinkResponse{
		Link: fullURL,
	}, nil
}

func (s *Service) PostCode(ctx context.Context, channelID, code string) error {
	if s.config.StreamlabsClientId == "" || s.config.StreamlabsClientSecret == "" {
		return errors.New("streamlabs integration not properly configured")
	}

	foundIntegration, err := s.streamlabsRepository.GetByChannelID(ctx, channelID)
	if err != nil && !errors.Is(err, streamlabs_integration.ErrNotFound) {
		return fmt.Errorf("failed to get streamlabs integration: %w", err)
	}

	redirectURL, err := s.getCallbackUrl(ctx)
	if err != nil {
		return fmt.Errorf("failed to get redirect URL: %w", err)
	}

	tokens, profile, err := s.getProfileData(
		ctx,
		s.config.StreamlabsClientId,
		s.config.StreamlabsClientSecret,
		redirectURL,
		code,
	)
	if err != nil {
		return fmt.Errorf("failed to get streamlabs profile data: %w", err)
	}

	if foundIntegration == model.Nil {
		if err := s.streamlabsRepository.Create(
			ctx, streamlabs_integration.CreateOpts{
				ChannelID:    channelID,
				AccessToken:  tokens.AccessToken,
				RefreshToken: tokens.RefreshToken,
				Enabled:      true,
				UserName:     profile.StreamLabs.DisplayName,
				Avatar:       profile.StreamLabs.ThumbNail,
			},
		); err != nil {
			return fmt.Errorf("failed to create streamlabs integration: %w", err)
		}
	} else {
		if err := s.streamlabsRepository.Update(
			ctx,
			streamlabs_integration.UpdateOpts{
				ChannelID:    channelID,
				AccessToken:  &tokens.AccessToken,
				RefreshToken: &tokens.RefreshToken,
				Enabled:      lo.ToPtr(true),
				UserName:     &profile.StreamLabs.DisplayName,
				Avatar:       &profile.StreamLabs.ThumbNail,
			},
		); err != nil {
			return fmt.Errorf("failed to update streamlabs integration: %w", err)
		}
	}

	newIntegration, err := s.streamlabsRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to get streamlabs integration after update: %w", err)
	}

	if err = s.twirBus.Integrations.Add.Publish(
		ctx, integrations.Request{
			ID:      newIntegration.ID.String(),
			Service: integrations.Streamlabs,
		},
	); err != nil {
		return fmt.Errorf("failed to publish add integration event: %w", err)
	}

	return err
}

func (s *Service) Logout(ctx context.Context, channelID string) error {
	err := s.streamlabsRepository.Delete(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to disable streamlabs integration: %w", err)
	}

	if err := s.twirBus.Integrations.Remove.Publish(
		ctx,
		integrations.Request{
			ID:      channelID,
			Service: integrations.Streamlabs,
		},
	); err != nil {
		return fmt.Errorf("failed to publish remove integration event: %w", err)
	}

	return nil
}

type streamlabsTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type streamlabsProfileResponse struct {
	StreamLabs struct {
		DisplayName string `json:"display_name"`
		ThumbNail   string `json:"thumbnail"`
		ID          int    `json:"id"`
	} `json:"streamlabs"`
}

func (s *Service) getProfileData(
	ctx context.Context,
	clientId, clientSecret, redirectURL, code string,
) (
	*streamlabsTokensResponse,
	*streamlabsProfileResponse,
	error,
) {
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", clientId)
	formData.Set("client_secret", clientSecret)
	formData.Set("redirect_uri", redirectURL)
	formData.Set("code", code)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://streamlabs.com/api/v2.0/token", bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to exchange code for tokens: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("failed to exchange code for tokens: %s", string(body))
	}

	var data streamlabsTokensResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	profileReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://streamlabs.com/api/v2.0/user", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create profile request: %w", err)
	}
	profileReq.Header.Set("Authorization", "Bearer "+data.AccessToken)

	profileResp, err := http.DefaultClient.Do(profileReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch streamlabs profile: %w", err)
	}
	defer profileResp.Body.Close()

	profileBody, err := io.ReadAll(profileResp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read profile response: %w", err)
	}

	if profileResp.StatusCode < 200 || profileResp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("failed to fetch streamlabs profile: %s", string(profileBody))
	}

	var profile streamlabsProfileResponse
	if err := json.Unmarshal(profileBody, &profile); err != nil {
		return nil, nil, fmt.Errorf("failed to parse profile response: %w", err)
	}

	return &data, &profile, nil
}
