package faceitintegration

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
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/integrations"
	config "github.com/twirapp/twir/libs/config"
	faceitintegrationentity "github.com/twirapp/twir/libs/entities/faceit_integration"
	faceitintegration "github.com/twirapp/twir/libs/repositories/faceit_integration"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	FaceitRepository faceitintegration.Repository
	TwirBus          *buscore.Bus
	Config           config.Config
	KV               kv.KV
}

func New(opts Opts) *Service {
	return &Service{
		faceitRepository: opts.FaceitRepository,
		twirBus:          opts.TwirBus,
		config:           opts.Config,
		kv:               opts.KV,
	}
}

type Service struct {
	faceitRepository faceitintegration.Repository
	twirBus          *buscore.Bus
	config           config.Config
	kv               kv.KV
}

type AuthLinkResponse struct {
	Link string `json:"link"`
}

func (s *Service) GetIntegrationData(ctx context.Context, channelID string) (
	faceitintegrationentity.Entity,
	error,
) {
	integration, err := s.faceitRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, faceitintegration.ErrNotFound) {
			return faceitintegrationentity.Entity{
				ChannelID: channelID,
				Enabled:   false,
			}, nil
		}
		return faceitintegrationentity.Entity{}, fmt.Errorf(
			"failed to get faceit integration: %w",
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

	return u.JoinPath("dashboard", "integrations", "faceit").String(), nil
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
	verifierValuer := s.kv.Get(ctx, fmt.Sprintf("faceit_pkce_%s", dashboardID))
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
		fmt.Sprintf("faceit_pkce_%s", dashboardID),
		verifier,
		kvoptions.WithExpire(15*time.Minute),
	); err != nil {
		return "", fmt.Errorf("failed to store PKCE code verifier: %w", err)
	}

	return verifier, nil
}

func (s *Service) GetAuthLink(ctx context.Context, dashboardID string) (*AuthLinkResponse, error) {
	if s.config.FaceitClientId == "" || s.config.FaceitClientSecret == "" {
		return nil, errors.New("faceit integration not properly configured")
	}

	codeVerifier, err := s.generatePkceCodeVerifier(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PKCE code verifier: %w", err)
	}

	codeChallange := generatePkceCodeChallenge(codeVerifier)

	redirectUrl, err := s.getCallbackUrl()
	if err != nil {
		return nil, fmt.Errorf("failed to get redirect URL: %w", err)
	}

	authURL := "https://accounts.faceit.com"
	params := url.Values{}
	params.Add("client_id", s.config.FaceitClientId)
	params.Add("redirect_popup", "true")
	params.Add("response_type", "code")
	params.Add("redirect_uri", redirectUrl)
	params.Add("code_challenge", codeChallange)
	params.Add("code_challenge_method", "S256")

	fullURL := fmt.Sprintf("%s?%s", authURL, params.Encode())

	return &AuthLinkResponse{
		Link: fullURL,
	}, nil
}

func (s *Service) PostCode(ctx context.Context, dashboardId, code string) error {
	if s.config.FaceitClientId == "" || s.config.FaceitClientSecret == "" {
		return errors.New("faceit integration not properly configured")
	}

	verifier, err := s.getPkceCodeVerifier(ctx, dashboardId)
	if err != nil {
		return fmt.Errorf("failed to get PKCE code verifier: %w", err)
	}

	foundIntegration, err := s.faceitRepository.GetByChannelID(ctx, dashboardId)
	if err != nil && !errors.Is(err, faceitintegration.ErrNotFound) {
		return fmt.Errorf("failed to get faceit integration: %w", err)
	}

	redirectUrl, err := s.getCallbackUrl()
	if err != nil {
		return fmt.Errorf("failed to get redirect URL: %w", err)
	}

	tokens, profile, err := s.getProfileData(
		ctx,
		s.config.FaceitClientId,
		s.config.FaceitClientSecret,
		redirectUrl,
		verifier,
		code,
	)
	if err != nil {
		return fmt.Errorf("failed to get faceit profile data: %w", err)
	}

	if foundIntegration.IsNil() {
		if err := s.faceitRepository.Create(
			ctx, faceitintegration.CreateOpts{
				ChannelID:   dashboardId,
				AccessToken: tokens.AccessToken,
				Enabled:     true,
				UserName:    profile.Nickname,
				Avatar:      profile.Avatar,
				Game:        "",
			},
		); err != nil {
			return fmt.Errorf("failed to create faceit integration: %w", err)
		}
	} else {
		if err := s.faceitRepository.Update(
			ctx,
			faceitintegration.UpdateOpts{
				ChannelID:   dashboardId,
				AccessToken: &tokens.AccessToken,
				Enabled:     lo.ToPtr(true),
				UserName:    &profile.Nickname,
				Avatar:      &profile.Avatar,
			},
		); err != nil {
			return fmt.Errorf("failed to update faceit integration: %w", err)
		}
	}

	newIntegration, err := s.faceitRepository.GetByChannelID(ctx, dashboardId)
	if err != nil {
		return fmt.Errorf("failed to get faceit integration after update: %w", err)
	}

	if err = s.twirBus.Integrations.Add.Publish(
		ctx, integrations.Request{
			ID:      fmt.Sprint(newIntegration.ID),
			Service: integrations.Faceit,
		},
	); err != nil {
		return fmt.Errorf("failed to publish add integration event: %w", err)
	}

	return nil
}

func (s *Service) UpdateGame(ctx context.Context, channelID, game string) error {
	err := s.faceitRepository.Update(
		ctx,
		faceitintegration.UpdateOpts{
			ChannelID: channelID,
			Game:      &game,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update faceit game: %w", err)
	}

	return nil
}

func (s *Service) Logout(ctx context.Context, channelID string) error {
	err := s.faceitRepository.Delete(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to disable faceit integration: %w", err)
	}

	if err := s.twirBus.Integrations.Remove.Publish(
		ctx,
		integrations.Request{
			ID:      channelID,
			Service: integrations.Faceit,
		},
	); err != nil {
		return fmt.Errorf("failed to publish remove integration event: %w", err)
	}

	return nil
}

type faceitTokensResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type faceitProfileResponse struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Country  string `json:"country"`
	GUID     string `json:"guid"`
}

func (s *Service) getProfileData(
	ctx context.Context,
	clientId, clientSecret, redirectURL, verifier, code string,
) (
	*faceitTokensResponse,
	*faceitProfileResponse,
	error,
) {
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", clientId)
	formData.Set("redirect_uri", redirectURL)
	formData.Set("code_verifier", verifier)
	formData.Set("code", code)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.faceit.com/auth/v1/oauth/token",
		bytes.NewBufferString(formData.Encode()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	authHeader := "Basic " + base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", clientId, clientSecret)),
	)
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read token response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf(
			"token request failed with status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	var tokens faceitTokensResponse
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	profileReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.faceit.com/auth/v1/resources/userinfo",
		nil,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create profile request: %w", err)
	}

	profileReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))

	profileResp, err := client.Do(profileReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute profile request: %w", err)
	}
	defer profileResp.Body.Close()

	profileBody, err := io.ReadAll(profileResp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read profile response body: %w", err)
	}

	if profileResp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf(
			"profile request failed with status %d: %s",
			profileResp.StatusCode,
			string(profileBody),
		)
	}

	var profile faceitProfileResponse
	if err := json.Unmarshal(profileBody, &profile); err != nil {
		return nil, nil, fmt.Errorf("failed to parse profile response: %w", err)
	}

	return &tokens, &profile, nil
}
