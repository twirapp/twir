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
	defaultIDAPIBaseURL  = "https://id.vk.ru"
	defaultIDHTTPTimeout = 10 * time.Second
	maxIDResponseSize    = int64(1 << 20)
)

var ErrDeviceIDRequired = errors.New("VK ID device ID is required")

type IDClientOpts struct {
	ClientID     string
	ServiceToken string
	RedirectURL  string
	APIBaseURL   string
	HTTPClient   *http.Client
}

type IDClient struct {
	clientID     string
	serviceToken string
	redirectURL  string
	apiBaseURL   string
	httpClient   *http.Client
}

type IDExchangeCodeInput struct {
	Code         string
	CodeVerifier string
	DeviceID     string
}

type IDRefreshTokenInput struct {
	RefreshToken string
	DeviceID     string
}

type IDToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	Scopes       []string
}

type IDUser struct {
	ID        string
	FirstName string
	LastName  string
	Avatar    string
}

type ProviderError struct {
	Code        string
	Description string
	State       string
	StatusCode  int
}

func (e *ProviderError) Error() string {
	if e.Code != "" && e.Description != "" {
		return fmt.Sprintf("VK ID error %s: %s", e.Code, e.Description)
	}
	if e.Code != "" {
		return fmt.Sprintf("VK ID error %s", e.Code)
	}
	if e.StatusCode != 0 {
		return fmt.Sprintf("VK ID returned HTTP status %d", e.StatusCode)
	}
	return "VK ID returned an unknown error"
}

func NewIDClient(opts IDClientOpts) (*IDClient, error) {
	if strings.TrimSpace(opts.ClientID) == "" {
		return nil, errors.New("VK ID client ID is required")
	}
	if strings.TrimSpace(opts.RedirectURL) == "" {
		return nil, errors.New("VK ID redirect URL is required")
	}
	if strings.TrimSpace(opts.ServiceToken) == "" {
		return nil, errors.New("VK ID service token is required")
	}

	apiBaseURL := opts.APIBaseURL
	if apiBaseURL == "" {
		apiBaseURL = defaultIDAPIBaseURL
	}
	parsedBaseURL, err := url.Parse(apiBaseURL)
	if err != nil || parsedBaseURL.Scheme == "" || parsedBaseURL.Host == "" {
		return nil, fmt.Errorf("invalid VK ID API base URL")
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultIDHTTPTimeout}
	}

	return &IDClient{
		clientID:     opts.ClientID,
		serviceToken: opts.ServiceToken,
		redirectURL:  opts.RedirectURL,
		apiBaseURL:   strings.TrimRight(apiBaseURL, "/"),
		httpClient:   httpClient,
	}, nil
}

func (c *IDClient) AuthorizationURL(state, codeChallenge string) (string, error) {
	endpoint, err := c.endpoint("authorize")
	if err != nil {
		return "", err
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parse VK ID authorization URL: %w", err)
	}
	query := u.Query()
	query.Set("response_type", "code")
	query.Set("client_id", c.clientID)
	query.Set("redirect_uri", c.redirectURL)
	query.Set("state", state)
	query.Set("code_challenge", codeChallenge)
	query.Set("code_challenge_method", "S256")
	query.Set("scope", "vkid.personal_info")
	u.RawQuery = query.Encode()

	return u.String(), nil
}

func (c *IDClient) ExchangeCode(ctx context.Context, input IDExchangeCodeInput) (*IDToken, error) {
	if strings.TrimSpace(input.DeviceID) == "" {
		return nil, ErrDeviceIDRequired
	}

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", c.clientID)
	form.Set("code", input.Code)
	form.Set("code_verifier", input.CodeVerifier)
	form.Set("redirect_uri", c.redirectURL)
	form.Set("device_id", input.DeviceID)
	c.addServiceToken(form)

	return c.exchangeToken(ctx, form)
}

func (c *IDClient) RefreshToken(ctx context.Context, input IDRefreshTokenInput) (*IDToken, error) {
	if strings.TrimSpace(input.DeviceID) == "" {
		return nil, ErrDeviceIDRequired
	}

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", c.clientID)
	form.Set("refresh_token", input.RefreshToken)
	form.Set("device_id", input.DeviceID)
	c.addServiceToken(form)

	return c.exchangeToken(ctx, form)
}

func (c *IDClient) UserInfo(ctx context.Context, accessToken string) (*IDUser, error) {
	form := url.Values{}
	form.Set("client_id", c.clientID)
	form.Set("access_token", accessToken)

	var response struct {
		User *struct {
			UserID    idUserID `json:"user_id"`
			FirstName string   `json:"first_name"`
			LastName  string   `json:"last_name"`
			Avatar    string   `json:"avatar"`
		} `json:"user"`
	}
	if err := c.postForm(ctx, "oauth2/user_info", form, &response); err != nil {
		return nil, err
	}
	if response.User == nil || response.User.UserID == "" {
		return nil, errors.New("VK ID user info response is missing user ID")
	}

	return &IDUser{
		ID:        string(response.User.UserID),
		FirstName: response.User.FirstName,
		LastName:  response.User.LastName,
		Avatar:    response.User.Avatar,
	}, nil
}

func (c *IDClient) exchangeToken(ctx context.Context, form url.Values) (*IDToken, error) {
	var response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
	}
	if err := c.postForm(ctx, "oauth2/auth", form, &response); err != nil {
		return nil, err
	}
	if strings.TrimSpace(response.AccessToken) == "" {
		return nil, errors.New("VK ID token response is missing access token")
	}

	return &IDToken{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		ExpiresIn:    response.ExpiresIn,
		Scopes:       strings.Fields(response.Scope),
	}, nil
}

func (c *IDClient) addServiceToken(form url.Values) {
	if c.serviceToken != "" {
		form.Set("service_token", c.serviceToken)
	}
}

func (c *IDClient) postForm(ctx context.Context, path string, form url.Values, target any) error {
	endpoint, err := c.endpoint(path)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return fmt.Errorf("create VK ID request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform VK ID request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxIDResponseSize+1))
	if err != nil {
		return fmt.Errorf("read VK ID response: %w", err)
	}
	if int64(len(body)) > maxIDResponseSize {
		return fmt.Errorf("VK ID response exceeds %d bytes", maxIDResponseSize)
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
		return fmt.Errorf("decode VK ID response: %w", err)
	}

	return nil
}

func (c *IDClient) endpoint(path string) (string, error) {
	endpoint, err := url.JoinPath(c.apiBaseURL, path)
	if err != nil {
		return "", fmt.Errorf("build VK ID endpoint: %w", err)
	}

	return endpoint, nil
}

type idUserID string

func (id *idUserID) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err == nil {
		*id = idUserID(value)
		return nil
	}

	var number json.Number
	if err := json.Unmarshal(data, &number); err != nil {
		return fmt.Errorf("decode VK ID user ID: %w", err)
	}
	*id = idUserID(number.String())

	return nil
}
