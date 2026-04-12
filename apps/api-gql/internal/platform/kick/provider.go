package kickplatform

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/twirapp/twir/apps/api-gql/internal/platform"
	cfg "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"go.uber.org/fx"
)

const (
	oauthBaseURL = "https://id.kick.com"
	apiBaseURL   = "https://api.kick.com"
)

var scopes = []string{"user:read", "events:subscribe", "chat:write", "channel:read"}

type Opts struct {
	fx.In

	Config cfg.Config
}

type Provider struct {
	config     cfg.Config
	httpClient *http.Client
}

func New(opts Opts) *Provider {
	return &Provider{
		config:     opts.Config,
		httpClient: &http.Client{},
	}
}

func (p *Provider) Name() string {
	return platformentity.PlatformKick.String()
}

func (p *Provider) GetAuthURL(state, codeChallenge string) string {
	query := url.Values{}
	query.Set("client_id", p.config.KickClientId)
	query.Set("redirect_uri", p.config.GetKickCallbackUrl())
	query.Set("response_type", "code")
	query.Set("scope", strings.Join(scopes, " "))
	query.Set("state", state)
	query.Set("code_challenge", codeChallenge)
	query.Set("code_challenge_method", "S256")

	return oauthBaseURL + "/oauth/authorize?" + query.Encode()
}

func (p *Provider) ExchangeCode(ctx context.Context, code, codeVerifier string) (*platform.PlatformTokens, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", p.config.KickClientId)
	form.Set("client_secret", p.config.KickClientSecret)
	form.Set("code", code)
	form.Set("code_verifier", codeVerifier)
	form.Set("redirect_uri", p.config.GetKickCallbackUrl())

	return p.requestTokens(ctx, form)
}

func (p *Provider) RefreshToken(ctx context.Context, refreshToken string) (*platform.PlatformTokens, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", p.config.KickClientId)
	form.Set("client_secret", p.config.KickClientSecret)
	form.Set("refresh_token", refreshToken)

	return p.requestTokens(ctx, form)
}

func (p *Provider) GetUser(ctx context.Context, accessToken string) (*platform.PlatformUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiBaseURL+"/public/v1/users", nil)
	if err != nil {
		return nil, fmt.Errorf("create kick user request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do kick user request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read kick user response: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("kick get user failed with status %d: %s", resp.StatusCode, string(body))
	}

	var parsed struct {
		Data []struct {
			UserID         int    `json:"user_id"`
			Name           string `json:"name"`
			Email          string `json:"email"`
			ProfilePicture string `json:"profile_picture"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("decode kick user response: %w", err)
	}

	if len(parsed.Data) == 0 {
		return nil, fmt.Errorf("kick user not found")
	}

	user := parsed.Data[0]

	return &platform.PlatformUser{
		ID:          strconv.Itoa(user.UserID),
		Login:       user.Name,
		DisplayName: user.Name,
		Avatar:      user.ProfilePicture,
	}, nil
}

func (p *Provider) requestTokens(ctx context.Context, form url.Values) (*platform.PlatformTokens, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		oauthBaseURL+"/oauth/token",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("create kick token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do kick token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read kick token response: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("kick token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var parsed struct {
		AccessToken  string   `json:"access_token"`
		RefreshToken string   `json:"refresh_token"`
		ExpiresIn    int      `json:"expires_in"`
		Scope        string   `json:"scope"`
		Scopes       []string `json:"scopes"`
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("decode kick token response: %w", err)
	}

	responseScopes := parsed.Scopes
	if len(responseScopes) == 0 && parsed.Scope != "" {
		responseScopes = strings.Fields(parsed.Scope)
	}

	return &platform.PlatformTokens{
		AccessToken:  parsed.AccessToken,
		RefreshToken: parsed.RefreshToken,
		ExpiresIn:    parsed.ExpiresIn,
		Scopes:       responseScopes,
	}, nil
}
