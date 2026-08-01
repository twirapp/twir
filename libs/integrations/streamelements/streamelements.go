package streamelements

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/twirapp/twir/libs/integrations/oauthlock"
)

const (
	defaultBaseURL   = "https://api.streamelements.com"
	providerName     = "streamelements"
	maxResponseBytes = int64(1 << 20)
)

var ErrUnauthorized = errors.New("streamelements unauthorized")

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type TokenStore interface {
	GetTokens(ctx context.Context, channelID string) (Tokens, error)
	UpdateTokens(ctx context.Context, channelID string, tokens Tokens) error
}

type Option func(*StreamElements)

func WithHTTPClient(client *http.Client) Option {
	return func(s *StreamElements) {
		if client != nil {
			s.httpClient = client
		}
	}
}

func WithBaseURL(baseURL string) Option {
	return func(s *StreamElements) {
		s.baseURL = strings.TrimRight(baseURL, "/")
	}
}

type StreamElements struct {
	clientID     string
	clientSecret string
	channelID    string
	redirectURL  string
	store        TokenStore
	locker       oauthlock.Locker
	httpClient   *http.Client
	baseURL      string

	tokensMu sync.RWMutex
	tokens   Tokens
}

func New(clientID, clientSecret string) *StreamElements {
	return NewStatic(clientID, clientSecret)
}

func NewStatic(clientID, clientSecret string, opts ...Option) *StreamElements {
	client := &StreamElements{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 15 * time.Second},
		baseURL:      defaultBaseURL,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func NewAuthorized(
	clientID, clientSecret, channelID, redirectURL string,
	tokens Tokens,
	store TokenStore,
	locker oauthlock.Locker,
	opts ...Option,
) *StreamElements {
	client := NewStatic(clientID, clientSecret, opts...)
	client.channelID = channelID
	client.redirectURL = redirectURL
	client.tokens = tokens
	client.store = store
	client.locker = locker
	return client
}

func (s *StreamElements) GetAuthLink(redirectURL string) string {
	return s.authLink(redirectURL, "")
}

func (s *StreamElements) GetAuthLinkWithState(redirectURL, state string) (string, error) {
	if strings.TrimSpace(state) == "" {
		return "", errors.New("StreamElements OAuth state is required")
	}
	return s.authLink(redirectURL, state), nil
}

func (s *StreamElements) authLink(redirectURL, state string) string {
	u, _ := url.Parse(s.baseURL + "/oauth2/authorize")
	query := u.Query()
	query.Set("client_id", s.clientID)
	query.Set("redirect_uri", redirectURL)
	query.Set("response_type", "code")
	query.Set("scope", "channel:read bot:read tips:read")
	if state != "" {
		query.Set("state", state)
	}
	u.RawQuery = query.Encode()
	return u.String()
}

func (s *StreamElements) ExchangeCode(
	ctx context.Context,
	code, redirectURL string,
) (*TokenResponse, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {s.clientID},
		"client_secret": {s.clientSecret},
		"code":          {code},
		"redirect_uri":  {redirectURL},
	}

	tokens := &TokenResponse{}
	if err := s.postToken(ctx, form, "exchange code", tokens); err != nil {
		return nil, err
	}
	s.setTokens(Tokens{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
	return tokens, nil
}

func (s *StreamElements) GetProfile(ctx context.Context) (*UserProfile, error) {
	profile := &UserProfile{}
	if err := s.authorizedJSON(ctx, http.MethodGet, "/kappa/v2/channels/me", profile); err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *StreamElements) GetCommands(ctx context.Context, channelID string) ([]Command, error) {
	commands := make([]Command, 0)
	path := fmt.Sprintf("/kappa/v2/bot/commands/%s", url.PathEscape(channelID))
	if err := s.authorizedJSON(ctx, http.MethodGet, path, &commands); err != nil {
		return nil, err
	}
	return commands, nil
}

func (s *StreamElements) GetTimers(ctx context.Context, channelID string) ([]Timer, error) {
	timers := make([]Timer, 0)
	path := fmt.Sprintf("/kappa/v2/bot/timers/%s", url.PathEscape(channelID))
	if err := s.authorizedJSON(ctx, http.MethodGet, path, &timers); err != nil {
		return nil, err
	}
	return timers, nil
}

func (s *StreamElements) authorizedJSON(
	ctx context.Context,
	method, path string,
	target any,
) error {
	current := s.currentTokens()
	request := func(accessToken string) error {
		return s.requestJSON(ctx, method, path, accessToken, target)
	}

	err := request(current.AccessToken)
	if !errors.Is(err, ErrUnauthorized) {
		return err
	}
	if s.store == nil || s.locker == nil {
		return ErrUnauthorized
	}

	err = s.locker.WithLock(ctx, s.lockKey(), func(ctx context.Context) error {
		fresh, err := s.store.GetTokens(ctx, s.channelID)
		if err != nil {
			return fmt.Errorf("reread StreamElements tokens: %w", err)
		}
		if fresh.AccessToken != current.AccessToken {
			current = fresh
			s.setTokens(fresh)
			return nil
		}

		rotated, err := s.refresh(ctx, fresh.RefreshToken)
		if err != nil {
			return err
		}
		if rotated.RefreshToken == "" {
			rotated.RefreshToken = fresh.RefreshToken
		}
		if err := s.store.UpdateTokens(ctx, s.channelID, rotated); err != nil {
			return fmt.Errorf("persist StreamElements tokens: %w", err)
		}
		current = rotated
		s.setTokens(rotated)
		return nil
	})
	if err != nil {
		return err
	}

	return request(current.AccessToken)
}

func (s *StreamElements) refresh(ctx context.Context, refreshToken string) (Tokens, error) {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {s.clientID},
		"client_secret": {s.clientSecret},
		"refresh_token": {refreshToken},
	}
	response := &TokenResponse{}
	if err := s.postToken(ctx, form, "refresh token", response); err != nil {
		return Tokens{}, err
	}
	return Tokens{AccessToken: response.AccessToken, RefreshToken: response.RefreshToken}, nil
}

func (s *StreamElements) postToken(
	ctx context.Context,
	form url.Values,
	operation string,
	target any,
) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.baseURL+"/oauth2/token",
		bytes.NewBufferString(form.Encode()),
	)
	if err != nil {
		return fmt.Errorf("create StreamElements %s request: %w", operation, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return s.doJSON(req, operation, target)
}

func (s *StreamElements) requestJSON(
	ctx context.Context,
	method, path, accessToken string,
	target any,
) error {
	req, err := http.NewRequestWithContext(ctx, method, s.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("create StreamElements API request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return s.doJSON(req, "API request", target)
}

func (s *StreamElements) doJSON(req *http.Request, operation string, target any) error {
	response, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform StreamElements %s: %w", operation, err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("StreamElements %s failed with status %d", operation, response.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil {
		return fmt.Errorf("read StreamElements %s response: %w", operation, err)
	}
	if int64(len(body)) > maxResponseBytes {
		return fmt.Errorf(
			"StreamElements %s response exceeds %d bytes",
			operation,
			maxResponseBytes,
		)
	}
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("decode StreamElements %s response: %w", operation, err)
	}
	return nil
}

func (s *StreamElements) currentTokens() Tokens {
	s.tokensMu.RLock()
	defer s.tokensMu.RUnlock()
	return s.tokens
}

func (s *StreamElements) setTokens(tokens Tokens) {
	s.tokensMu.Lock()
	defer s.tokensMu.Unlock()
	s.tokens = tokens
}

func (s *StreamElements) lockKey() string {
	return "twir:integration-token-refresh:" + providerName + ":" + s.channelID
}
