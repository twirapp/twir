package vk_integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/samber/lo"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/vk_integration"
	vkintegrationrepo "github.com/twirapp/twir/libs/repositories/vk_integration"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	VKRepository vkintegrationrepo.Repository
	Config       config.Config
}

func New(opts Opts) *Service {
	return &Service{
		vkRepository: opts.VKRepository,
		config:       opts.Config,
	}
}

type Service struct {
	vkRepository vkintegrationrepo.Repository
	config       config.Config
}

type AuthLinkResponse struct {
	Link string `json:"link"`
}

type vkProfile struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PhotoMaxOrig string `json:"photo_max_orig"`
	ID           int    `json:"id"`
}

type vkProfileResponse struct {
	Error *struct {
		Msg  string `json:"error_msg"`
		Code int    `json:"error_code"`
	}
	Response []vkProfile `json:"response"`
}

type vkTokensResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *Service) GetIntegrationData(ctx context.Context, channelID string) (
	vk_integration.Entity,
	error,
) {
	integration, err := s.vkRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, vkintegrationrepo.ErrNotFound) {
			// Return default data if integration doesn't exist
			return vk_integration.Entity{
				ChannelID: channelID,
				Enabled:   false,
			}, nil
		}
		return vk_integration.Entity{}, fmt.Errorf(
			"failed to get vk integration: %w",
			err,
		)
	}

	return integration, nil
}

func (s *Service) getCallbackUrl() (string, error) {
	u, err := url.Parse(s.config.SiteBaseUrl)
	if err != nil {
		return "", fmt.Errorf("invalid site base URL: %w", err)
	}

	return u.JoinPath("dashboard", "integrations", "vk").String(), nil
}

func (s *Service) GetAuthLink(ctx context.Context) (*AuthLinkResponse, error) {
	if s.config.VKClientId == "" || s.config.VKClientSecret == "" {
		return nil, errors.New("vk integration not properly configured")
	}

	redirectUrl, err := s.getCallbackUrl()
	if err != nil {
		return nil, fmt.Errorf("failed to get redirect URL: %w", err)
	}

	// Build OAuth authorization URL
	link, _ := url.Parse("https://oauth.vk.com/authorize")
	q := link.Query()
	q.Add("client_id", s.config.VKClientId)
	q.Add("display", "page")
	q.Add("response_type", "code")
	q.Add("scope", "status offline")
	q.Add("redirect_uri", redirectUrl)
	link.RawQuery = q.Encode()

	return &AuthLinkResponse{
		Link: link.String(),
	}, nil
}

func (s *Service) getProfileData(
	ctx context.Context,
	clientID string,
	clientSecret string,
	redirectUrl string,
	code string,
) (*vkTokensResponse, *vkProfileResponse, error) {
	tokenUrl, _ := url.Parse("https://oauth.vk.com/access_token")
	q := tokenUrl.Query()
	q.Set("grant_type", "authorization_code")
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("redirect_uri", redirectUrl)
	q.Set("code", code)
	tokenUrl.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenUrl.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("vk auth error: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	tokensData := &vkTokensResponse{}
	if err := json.Unmarshal(bodyBytes, tokensData); err != nil {
		return nil, nil, err
	}

	profileUrl, _ := url.Parse("https://api.vk.com/method/users.get")
	pq := profileUrl.Query()
	pq.Set("v", "5.131")
	pq.Set("fields", "photo_max_orig")
	pq.Set("access_token", tokensData.AccessToken)
	profileUrl.RawQuery = pq.Encode()

	profileReq, err := http.NewRequestWithContext(ctx, http.MethodGet, profileUrl.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	profileResp, err := http.DefaultClient.Do(profileReq)
	if err != nil {
		return nil, nil, err
	}
	defer profileResp.Body.Close()

	if profileResp.StatusCode < 200 || profileResp.StatusCode >= 300 {
		profileBodyBytes, _ := io.ReadAll(profileResp.Body)
		return nil, nil, fmt.Errorf("vk profile error: %s", string(profileBodyBytes))
	}

	profileBodyBytes, err := io.ReadAll(profileResp.Body)
	if err != nil {
		return nil, nil, err
	}

	profileData := &vkProfileResponse{}
	if err := json.Unmarshal(profileBodyBytes, profileData); err != nil {
		return nil, nil, err
	}

	if profileData.Error != nil {
		return nil, nil, fmt.Errorf(
			"vk profile error: %s (code: %d)",
			profileData.Error.Msg,
			profileData.Error.Code,
		)
	}

	return tokensData, profileData, nil
}

func (s *Service) PostCode(ctx context.Context, channelID, code string) error {
	if s.config.VKClientId == "" || s.config.VKClientSecret == "" {
		return errors.New("vk integration not properly configured")
	}

	foundIntegration, err := s.vkRepository.GetByChannelID(ctx, channelID)
	if err != nil && !errors.Is(err, vkintegrationrepo.ErrNotFound) {
		return fmt.Errorf("failed to get vk integration: %w", err)
	}

	redirectUrl, err := s.getCallbackUrl()
	if err != nil {
		return fmt.Errorf("failed to get redirect URL: %w", err)
	}

	tokens, profile, err := s.getProfileData(
		ctx,
		s.config.VKClientId,
		s.config.VKClientSecret,
		redirectUrl,
		code,
	)
	if err != nil {
		return fmt.Errorf("failed to get vk profile data: %w", err)
	}

	userName := profile.Response[0].FirstName + " " + profile.Response[0].LastName

	if foundIntegration.IsNil() {
		if err := s.vkRepository.Create(
			ctx, vkintegrationrepo.CreateOpts{
				ChannelID:   channelID,
				AccessToken: tokens.AccessToken,
				Enabled:     true,
				UserName:    userName,
				Avatar:      profile.Response[0].PhotoMaxOrig,
			},
		); err != nil {
			return fmt.Errorf("failed to create vk integration: %w", err)
		}
	} else {
		if err := s.vkRepository.Update(
			ctx,
			vkintegrationrepo.UpdateOpts{
				ChannelID:   channelID,
				AccessToken: &tokens.AccessToken,
				Enabled:     lo.ToPtr(true),
				UserName:    &userName,
				Avatar:      &profile.Response[0].PhotoMaxOrig,
			},
		); err != nil {
			return fmt.Errorf("failed to update vk integration: %w", err)
		}
	}

	return nil
}

func (s *Service) Logout(ctx context.Context, channelID string) error {
	err := s.vkRepository.Delete(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to delete vk integration: %w", err)
	}

	return nil
}
