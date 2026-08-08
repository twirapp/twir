package donatex_integration

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/integrations"
	config "github.com/twirapp/twir/libs/config"
	donatexintegrationentity "github.com/twirapp/twir/libs/entities/donatex_integration"
	donatexintegration "github.com/twirapp/twir/libs/repositories/donatex_integration"
)

const (
	authorizeURL = "https://donatex.gg/api/connect/authorize"
	tokenURL     = "https://donatex.gg/api/connect/token"
	profileURL   = "https://donatex.gg/api/v1/user/me"
	scopes       = "openid offline_access user.read donations.read donations.subscribe"
)

func New(
	donateXRepository donatexintegration.Repository,
	twirBus *buscore.Bus,
	config config.Config,
	kvClient kv.KV,
) *Service {
	return &Service{
		donateXRepository: donateXRepository,
		twirBus:           twirBus,
		config:            config,
		kv:                kvClient,
	}
}

type Service struct {
	donateXRepository donatexintegration.Repository
	twirBus           *buscore.Bus
	config            config.Config
	kv                kv.KV
}

type AuthLinkResponse struct {
	Link string `json:"link"`
}

func (s *Service) GetIntegrationData(ctx context.Context, channelID string) (
	donatexintegrationentity.Entity,
	error,
) {
	integration, err := s.donateXRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return donatexintegrationentity.Entity{}, fmt.Errorf(
			"failed to get donatex integration: %w",
			err,
		)
	}

	if integration.IsNil() {
		return donatexintegrationentity.Entity{
			ChannelID: channelID,
			Enabled:   false,
		}, nil
	}

	return integration, nil
}

func (s *Service) getCallbackUrl(ctx context.Context) (string, error) {
	baseUrl, _ := gincontext.GetBaseUrlFromContext(ctx, s.config.SiteBaseUrl)
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("invalid site base URL: %w", err)
	}

	return u.JoinPath("dashboard", "integrations", "donatex").String(), nil
}

func generatePkceCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	b64 := base64.StdEncoding.EncodeToString(h[:])
	b64 = strings.ReplaceAll(b64, "+", "-")
	b64 = strings.ReplaceAll(b64, "/", "_")
	b64 = strings.TrimRight(b64, "=")
	return b64
}

func (s *Service) getPkceCodeVerifier(ctx context.Context, dashboardID string) (string, error) {
	verifierValuer := s.kv.Get(ctx, fmt.Sprintf("donatex_pkce_%s", dashboardID))
	if err := verifierValuer.Err(); err != nil {
		return "", fmt.Errorf("failed to get PKCE code verifier: %w", err)
	}

	verifier, err := verifierValuer.String()
	if err != nil {
		return "", fmt.Errorf("failed to parse PKCE code verifier: %w", err)
	}

	return verifier, nil
}

func (s *Service) generatePkceCodeVerifier(ctx context.Context, dashboardID string) (string, error) {
	b := make([]byte, 48)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	verifier := hex.EncodeToString(b)

	if err := s.kv.Set(
		ctx,
		fmt.Sprintf("donatex_pkce_%s", dashboardID),
		verifier,
		kvoptions.WithExpire(15*time.Minute),
	); err != nil {
		return "", fmt.Errorf("failed to store PKCE code verifier: %w", err)
	}

	return verifier, nil
}

func (s *Service) GetAuthLink(ctx context.Context, dashboardID string) (*AuthLinkResponse, error) {
	if s.config.DonateXClientId == "" || s.config.DonateXClientSecret == "" {
		return nil, errors.New("donatex integration not properly configured")
	}

	codeVerifier, err := s.generatePkceCodeVerifier(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PKCE code verifier: %w", err)
	}

	codeChallenge := generatePkceCodeChallenge(codeVerifier)

	redirectUrl, err := s.getCallbackUrl(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get redirect URL: %w", err)
	}

	params := url.Values{}
	params.Add("client_id", s.config.DonateXClientId)
	params.Add("redirect_uri", redirectUrl)
	params.Add("response_type", "code")
	params.Add("scope", scopes)
	params.Add("code_challenge", codeChallenge)
	params.Add("code_challenge_method", "S256")

	fullURL := fmt.Sprintf("%s?%s", authorizeURL, params.Encode())

	return &AuthLinkResponse{
		Link: fullURL,
	}, nil
}

func (s *Service) PostCode(ctx context.Context, channelID, code string) error {
	if s.config.DonateXClientId == "" || s.config.DonateXClientSecret == "" {
		return errors.New("donatex integration not properly configured")
	}

	verifier, err := s.getPkceCodeVerifier(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to get PKCE code verifier: %w", err)
	}

	foundIntegration, err := s.donateXRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to get donatex integration: %w", err)
	}

	redirectUrl, err := s.getCallbackUrl(ctx)
	if err != nil {
		return fmt.Errorf("failed to get redirect URL: %w", err)
	}

	tokens, profile, err := s.getProfileData(
		ctx,
		s.config.DonateXClientId,
		s.config.DonateXClientSecret,
		redirectUrl,
		verifier,
		code,
	)
	if err != nil {
		return fmt.Errorf("failed to get donatex profile data: %w", err)
	}

	if foundIntegration.IsNil() {
		if err := s.donateXRepository.Create(
			ctx, donatexintegration.CreateOpts{
				ChannelID:     channelID,
				AccessToken:   tokens.AccessToken,
				RefreshToken:  tokens.RefreshToken,
				DonateXUserID: profile.ID,
				Enabled:       true,
				UserName:      profile.UserName,
				Avatar:        profile.AvatarUrl,
			},
		); err != nil {
			return fmt.Errorf("failed to create donatex integration: %w", err)
		}
	} else {
		if err := s.donateXRepository.Update(
			ctx,
			donatexintegration.UpdateOpts{
				ChannelID:     channelID,
				AccessToken:   &tokens.AccessToken,
				RefreshToken:  &tokens.RefreshToken,
				DonateXUserID: &profile.ID,
				Enabled:       lo.ToPtr(true),
				UserName:      &profile.UserName,
				Avatar:        &profile.AvatarUrl,
			},
		); err != nil {
			return fmt.Errorf("failed to update donatex integration: %w", err)
		}
	}

	newIntegration, err := s.donateXRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to get donatex integration after update: %w", err)
	}

	if err = s.twirBus.Integrations.Add.Publish(
		ctx, integrations.Request{
			ID:      fmt.Sprint(newIntegration.ID),
			Service: integrations.DonateX,
		},
	); err != nil {
		return fmt.Errorf("failed to publish add integration event: %w", err)
	}

	return nil
}

func (s *Service) Logout(ctx context.Context, channelID string) error {
	if err := s.donateXRepository.Delete(ctx, channelID); err != nil {
		return fmt.Errorf("failed to disable donatex integration: %w", err)
	}

	if err := s.twirBus.Integrations.Remove.Publish(
		ctx,
		integrations.Request{
			ID:      channelID,
			Service: integrations.DonateX,
		},
	); err != nil {
		return fmt.Errorf("failed to publish remove integration event: %w", err)
	}

	return nil
}

type donatexTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type donatexProfileResponse struct {
	ID        string `json:"id"`
	UserName  string `json:"username"`
	AvatarUrl string `json:"avatarUrl"`
}

func (s *Service) getProfileData(
	ctx context.Context,
	clientId, clientSecret, redirectURL, verifier, code string,
) (
	*donatexTokensResponse,
	*donatexProfileResponse,
	error,
) {
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", clientId)
	formData.Set("client_secret", clientSecret)
	formData.Set("redirect_uri", redirectURL)
	formData.Set("code_verifier", verifier)
	formData.Set("code", code)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		tokenURL,
		bytes.NewBufferString(formData.Encode()),
	)
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

	var tokens donatexTokensResponse
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	profileReq, err := http.NewRequestWithContext(ctx, http.MethodGet, profileURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create profile request: %w", err)
	}
	profileReq.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

	profileResp, err := http.DefaultClient.Do(profileReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch donatex profile: %w", err)
	}
	defer profileResp.Body.Close()

	profileBody, err := io.ReadAll(profileResp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read profile response: %w", err)
	}

	if profileResp.StatusCode < 200 || profileResp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("failed to fetch donatex profile: %s", string(profileBody))
	}

	var profile donatexProfileResponse
	if err := json.Unmarshal(profileBody, &profile); err != nil {
		return nil, nil, fmt.Errorf("failed to parse profile response: %w", err)
	}

	return &tokens, &profile, nil
}
