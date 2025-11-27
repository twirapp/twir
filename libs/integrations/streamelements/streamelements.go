package streamelements

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	baseURL = "https://api.streamelements.com"
)

type StreamElements struct {
	clientID     string
	clientSecret string
	accessToken  string
}

func New(clientID, clientSecret string) *StreamElements {
	return &StreamElements{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

var authScopes = []string{
	"channel:read",
	"bot:read",
	"loyalty:read",
	"activities:read",
	"tips:read",
	"activities:read",
	"overlays:read",
}

// GetAuthLink generates the OAuth2 authorization URL
func (s *StreamElements) GetAuthLink(redirectURL string) string {
	u, _ := url.Parse(baseURL + "/oauth2/authorize")

	q := u.Query()
	q.Add("client_id", s.clientID)
	q.Add("redirect_uri", redirectURL)
	q.Add("response_type", "code")
	q.Add(
		"scope",
		strings.Join(authScopes, " "),
	)
	u.RawQuery = q.Encode()

	return u.String()
}

// ExchangeCode exchanges the authorization code for access token
func (s *StreamElements) ExchangeCode(
	ctx context.Context,
	code, redirectURL string,
) (*TokenResponse, error) {
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", s.clientID)
	formData.Set("client_secret", s.clientSecret)
	formData.Set("code", code)
	formData.Set("redirect_uri", redirectURL)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/oauth2/token",
		bytes.NewBufferString(formData.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to exchange code: %s", string(body))
	}

	tokenData := &TokenResponse{}
	if err := json.Unmarshal(body, tokenData); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	s.accessToken = tokenData.AccessToken
	return tokenData, nil
}

// GetProfile fetches the user profile
func (s *StreamElements) GetProfile(ctx context.Context) (*UserProfile, error) {
	if s.accessToken == "" {
		return nil, fmt.Errorf("no access token available")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/kappa/v2/channels/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to get profile: %s", string(body))
	}

	profile := &UserProfile{}
	if err := json.Unmarshal(body, profile); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return profile, nil
}

// GetCommands fetches all commands for a channel
func (s *StreamElements) GetCommands(ctx context.Context, channelID string) ([]Command, error) {
	if s.accessToken == "" {
		return nil, fmt.Errorf("no access token available")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/kappa/v2/bot/commands/%s", baseURL, channelID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to get commands: %s", string(body))
	}

	var commands []Command
	if err := json.Unmarshal(body, &commands); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return commands, nil
}

func (s *StreamElements) GetTimers(ctx context.Context, channelID string) ([]Timer, error) {
	if s.accessToken == "" {
		return nil, fmt.Errorf("no access token available")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/kappa/v2/bot/timers/%s", baseURL, channelID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to get timers: %s", string(body))
	}

	var timers []Timer
	if err := json.Unmarshal(body, &timers); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return timers, nil
}
