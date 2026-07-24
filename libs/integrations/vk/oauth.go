package vk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultOAuthAPIBaseURL    = "https://api.live.vkvideo.ru"
	defaultOAuthAuthBaseURL   = "https://auth.live.vkvideo.ru"
	defaultOAuthDevAPIBaseURL = "https://apidev.live.vkvideo.ru"
	defaultOAuthHTTPTimeout   = 10 * time.Second
	maxOAuthResponseSize      = int64(1 << 20)
)

type OAuthClientOpts struct {
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	APIBaseURL    string
	AuthBaseURL   string
	DevAPIBaseURL string
	HTTPClient    *http.Client
}

type OAuthClient struct {
	clientID      string
	clientSecret  string
	redirectURL   string
	apiBaseURL    string
	authBaseURL   string
	devAPIBaseURL string
	httpClient    *http.Client
}

type OAuthToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	Scopes       []string
}

type VideoUser struct {
	ID     string
	Nick   string
	Avatar string
}

type ProviderError struct {
	Code        string
	Description string
	State       string
	StatusCode  int
}

func (e *ProviderError) Error() string {
	if e.Code != "" && e.Description != "" {
		return fmt.Sprintf("VK Video error %s: %s", e.Code, e.Description)
	}
	if e.Code != "" {
		return fmt.Sprintf("VK Video error %s", e.Code)
	}
	if e.StatusCode != 0 {
		return fmt.Sprintf("VK Video returned HTTP status %d", e.StatusCode)
	}
	return "VK Video returned an unknown error"
}

func NewOAuthClient(opts OAuthClientOpts) (*OAuthClient, error) {
	if strings.TrimSpace(opts.ClientID) == "" {
		return nil, errors.New("VK Video client ID is required")
	}
	if strings.TrimSpace(opts.ClientSecret) == "" {
		return nil, errors.New("VK Video client secret is required")
	}
	if strings.TrimSpace(opts.RedirectURL) == "" {
		return nil, errors.New("VK Video redirect URL is required")
	}

	apiBaseURL := opts.APIBaseURL
	if apiBaseURL == "" {
		apiBaseURL = defaultOAuthAPIBaseURL
	}
	if err := validateBaseURL(apiBaseURL); err != nil {
		return nil, fmt.Errorf("invalid VK Video API base URL: %w", err)
	}

	authBaseURL := opts.AuthBaseURL
	if authBaseURL == "" {
		authBaseURL = defaultOAuthAuthBaseURL
	}
	if err := validateBaseURL(authBaseURL); err != nil {
		return nil, fmt.Errorf("invalid VK Video auth base URL: %w", err)
	}

	devAPIBaseURL := opts.DevAPIBaseURL
	if devAPIBaseURL == "" {
		devAPIBaseURL = defaultOAuthDevAPIBaseURL
	}
	if err := validateBaseURL(devAPIBaseURL); err != nil {
		return nil, fmt.Errorf("invalid VK Video DevAPI base URL: %w", err)
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultOAuthHTTPTimeout}
	}

	return &OAuthClient{
		clientID:      opts.ClientID,
		clientSecret:  opts.ClientSecret,
		redirectURL:   opts.RedirectURL,
		apiBaseURL:    strings.TrimRight(apiBaseURL, "/"),
		authBaseURL:   strings.TrimRight(authBaseURL, "/"),
		devAPIBaseURL: strings.TrimRight(devAPIBaseURL, "/"),
		httpClient:    httpClient,
	}, nil
}

func (c *OAuthClient) AuthorizationURL(state string, scopes []string) (string, error) {
	endpoint, err := url.JoinPath(c.authBaseURL, "app", "oauth2", "authorize")
	if err != nil {
		return "", fmt.Errorf("build VK Video authorization endpoint: %w", err)
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parse VK Video authorization URL: %w", err)
	}
	query := u.Query()
	query.Set("response_type", "code")
	query.Set("client_id", c.clientID)
	query.Set("redirect_uri", c.redirectURL)
	query.Set("state", state)
	if len(scopes) > 0 {
		query.Set("scope", strings.Join(scopes, ","))
	}
	u.RawQuery = query.Encode()

	return u.String(), nil
}

func (c *OAuthClient) ExchangeCode(ctx context.Context, code string) (*OAuthToken, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", c.redirectURL)

	return c.requestToken(ctx, form)
}

func (c *OAuthClient) RefreshToken(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)
	form.Set("redirect_uri", c.redirectURL)

	return c.requestToken(ctx, form)
}

func (c *OAuthClient) CurrentUser(ctx context.Context, accessToken string) (*VideoUser, error) {
	endpoint, err := url.JoinPath(c.devAPIBaseURL, "v1", "current_user")
	if err != nil {
		return nil, fmt.Errorf("build VK Video current user endpoint: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create VK Video request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	var response struct {
		Data struct {
			User *struct {
				ID        flexString `json:"id"`
				Nick      string     `json:"nick"`
				AvatarURL string     `json:"avatar_url"`
			} `json:"user"`
		} `json:"data"`
	}
	if err := c.do(req, &response); err != nil {
		return nil, err
	}
	if response.Data.User == nil || response.Data.User.ID == "" {
		return nil, errors.New("VK Video current user response is missing user ID")
	}

	return &VideoUser{
		ID:     string(response.Data.User.ID),
		Nick:   response.Data.User.Nick,
		Avatar: response.Data.User.AvatarURL,
	}, nil
}

func (c *OAuthClient) requestToken(ctx context.Context, form url.Values) (*OAuthToken, error) {
	endpoint, err := url.JoinPath(c.apiBaseURL, "oauth", "server", "token")
	if err != nil {
		return nil, fmt.Errorf("build VK Video token endpoint: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("create VK Video request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.clientID, c.clientSecret)

	var response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
	}
	if err := c.do(req, &response); err != nil {
		return nil, err
	}
	if strings.TrimSpace(response.AccessToken) == "" {
		return nil, errors.New("VK Video token response is missing access token")
	}

	return &OAuthToken{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		ExpiresIn:    response.ExpiresIn,
		Scopes:       strings.Fields(response.Scope),
	}, nil
}

func (c *OAuthClient) do(req *http.Request, target any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform VK Video request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxOAuthResponseSize+1))
	if err != nil {
		return fmt.Errorf("read VK Video response: %w", err)
	}
	if int64(len(body)) > maxOAuthResponseSize {
		return fmt.Errorf("VK Video response exceeds %d bytes", maxOAuthResponseSize)
	}

	var providerResponse struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
		State            string `json:"state"`
	}
	if err := json.Unmarshal(body, &providerResponse); err == nil && providerResponse.Error != "" {
		return &ProviderError{
			Code:        providerResponse.Error,
			Description: providerResponse.ErrorDescription,
			State:       providerResponse.State,
			StatusCode:  resp.StatusCode,
		}
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return &ProviderError{StatusCode: resp.StatusCode}
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("decode VK Video response: %w", err)
	}

	return nil
}

func validateBaseURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return errors.New("must be an absolute URL")
	}

	return nil
}

type flexString string

func (s *flexString) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err == nil {
		*s = flexString(value)
		return nil
	}

	var number json.Number
	if err := json.Unmarshal(data, &number); err != nil {
		return fmt.Errorf("decode VK Video identifier: %w", err)
	}
	*s = flexString(number.String())

	return nil
}
