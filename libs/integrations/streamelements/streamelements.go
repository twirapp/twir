package streamelements

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/imroc/req/v3"
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
	tokenData := &TokenResponse{}

	resp, err := req.R().
		SetContext(ctx).
		SetFormData(
			map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     s.clientID,
				"client_secret": s.clientSecret,
				"code":          code,
				"redirect_uri":  redirectURL,
			},
		).
		SetSuccessResult(tokenData).
		Post(baseURL + "/oauth2/token")

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to exchange code: %s", resp.String())
	}

	s.accessToken = tokenData.AccessToken
	return tokenData, nil
}

// GetProfile fetches the user profile
func (s *StreamElements) GetProfile(ctx context.Context) (*UserProfile, error) {
	if s.accessToken == "" {
		return nil, fmt.Errorf("no access token available")
	}

	profile := &UserProfile{}

	resp, err := req.R().
		SetContext(ctx).
		SetBearerAuthToken(s.accessToken).
		SetSuccessResult(profile).
		Get(baseURL + "/kappa/v2/channels/me")

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to get profile: %s", resp.String())
	}

	return profile, nil
}

// GetCommands fetches all commands for a channel
func (s *StreamElements) GetCommands(ctx context.Context, channelID string) ([]Command, error) {
	if s.accessToken == "" {
		return nil, fmt.Errorf("no access token available")
	}

	var commands []Command

	resp, err := req.R().
		SetContext(ctx).
		SetBearerAuthToken(s.accessToken).
		SetSuccessResult(&commands).
		Get(fmt.Sprintf("%s/kappa/v2/bot/commands/%s", baseURL, channelID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to get commands: %s", resp.String())
	}

	return commands, nil
}

func (s *StreamElements) GetTimers(ctx context.Context, channelID string) ([]Timer, error) {
	if s.accessToken == "" {
		return nil, fmt.Errorf("no access token available")
	}

	var timers []Timer

	resp, err := req.R().
		SetContext(ctx).
		SetBearerAuthToken(s.accessToken).
		SetSuccessResult(&timers).
		Get(fmt.Sprintf("%s/kappa/v2/bot/timers/%s", baseURL, channelID))

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to get timers: %s", resp.String())
	}

	return timers, nil
}
