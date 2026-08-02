package streamlabs

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
	defaultAPIBaseURL  = "https://streamlabs.com"
	defaultAuthBaseURL = "https://www.streamlabs.com"
	providerName       = "streamlabs"
	maxResponseBytes   = int64(1 << 20)
)

var ErrUnauthorized = errors.New("streamlabs unauthorized")

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type TokenStore interface {
	GetTokens(ctx context.Context, channelID string) (Tokens, error)
	UpdateTokens(ctx context.Context, channelID string, tokens Tokens) error
}

type Option func(*Streamlabs)

func WithHTTPClient(client *http.Client) Option {
	return func(s *Streamlabs) {
		if client != nil {
			s.httpClient = client
		}
	}
}

func WithBaseURL(baseURL string) Option {
	return func(s *Streamlabs) {
		baseURL = strings.TrimRight(baseURL, "/")
		s.apiBaseURL = baseURL
		s.authBaseURL = baseURL
	}
}

type Streamlabs struct {
	clientID     string
	clientSecret string
	channelID    string
	redirectURL  string
	store        TokenStore
	locker       oauthlock.Locker
	httpClient   *http.Client
	apiBaseURL   string
	authBaseURL  string

	tokensMu sync.RWMutex
	tokens   Tokens
}

func New(clientID, clientSecret, redirectURL string, opts ...Option) *Streamlabs {
	client := &Streamlabs{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		httpClient:   &http.Client{Timeout: 15 * time.Second},
		apiBaseURL:   defaultAPIBaseURL,
		authBaseURL:  defaultAuthBaseURL,
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
) *Streamlabs {
	client := New(clientID, clientSecret, redirectURL, opts...)
	client.channelID = channelID
	client.tokens = tokens
	client.store = store
	client.locker = locker
	return client
}

func (s *Streamlabs) GetAuthLink(state ...string) string {
	u, _ := url.Parse(s.authBaseURL + "/api/v2.0/authorize")
	query := u.Query()
	query.Set("client_id", s.clientID)
	query.Set("redirect_uri", s.redirectURL)
	query.Set("response_type", "code")
	query.Set("scope", "socket.token donations.read")
	if len(state) > 0 && state[0] != "" {
		query.Set("state", state[0])
	}
	u.RawQuery = query.Encode()
	return u.String()
}

func (s *Streamlabs) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {s.clientID},
		"client_secret": {s.clientSecret},
		"code":          {code},
		"redirect_uri":  {s.redirectURL},
	}

	tokens := &TokenResponse{}
	if err := s.postToken(ctx, form, "exchange code", tokens); err != nil {
		return nil, err
	}
	s.setTokens(Tokens{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
	return tokens, nil
}

func (s *Streamlabs) GetProfile(ctx context.Context) (*UserProfile, error) {
	profile := &UserProfile{}
	if err := s.authorizedJSON(ctx, http.MethodGet, "/api/v2.0/user", profile); err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *Streamlabs) GetSocketToken(ctx context.Context) (*SocketTokenResponse, error) {
	token := &SocketTokenResponse{}
	if err := s.authorizedJSON(ctx, http.MethodGet, "/api/v2.0/socket/token", token); err != nil {
		return nil, err
	}
	return token, nil
}

func (s *Streamlabs) authorizedJSON(
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
			return fmt.Errorf("reread Streamlabs tokens: %w", err)
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
			return fmt.Errorf("persist Streamlabs tokens: %w", err)
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

func (s *Streamlabs) refresh(ctx context.Context, refreshToken string) (Tokens, error) {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {s.clientID},
		"client_secret": {s.clientSecret},
		"refresh_token": {refreshToken},
		"redirect_uri":  {s.redirectURL},
	}
	response := &TokenResponse{}
	if err := s.postToken(ctx, form, "refresh token", response); err != nil {
		return Tokens{}, err
	}
	return Tokens{AccessToken: response.AccessToken, RefreshToken: response.RefreshToken}, nil
}

func (s *Streamlabs) postToken(
	ctx context.Context,
	form url.Values,
	operation string,
	target any,
) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.apiBaseURL+"/api/v2.0/token",
		bytes.NewBufferString(form.Encode()),
	)
	if err != nil {
		return fmt.Errorf("create Streamlabs %s request: %w", operation, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return s.doJSON(req, operation, target)
}

func (s *Streamlabs) requestJSON(
	ctx context.Context,
	method, path, accessToken string,
	target any,
) error {
	req, err := http.NewRequestWithContext(ctx, method, s.apiBaseURL+path, nil)
	if err != nil {
		return fmt.Errorf("create Streamlabs API request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	return s.doJSON(req, "API request", target)
}

func (s *Streamlabs) doJSON(req *http.Request, operation string, target any) error {
	response, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("perform Streamlabs %s: %w", operation, err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("Streamlabs %s failed with status %d", operation, response.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil {
		return fmt.Errorf("read Streamlabs %s response: %w", operation, err)
	}
	if int64(len(body)) > maxResponseBytes {
		return fmt.Errorf(
			"Streamlabs %s response exceeds %d bytes",
			operation,
			maxResponseBytes,
		)
	}
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("decode Streamlabs %s response: %w", operation, err)
	}
	return nil
}

func (s *Streamlabs) currentTokens() Tokens {
	s.tokensMu.RLock()
	defer s.tokensMu.RUnlock()
	return s.tokens
}

func (s *Streamlabs) setTokens(tokens Tokens) {
	s.tokensMu.Lock()
	defer s.tokensMu.Unlock()
	s.tokens = tokens
}

func (s *Streamlabs) lockKey() string {
	return RefreshLockKey(s.channelID)
}

func RefreshLockKey(channelID string) string {
	return "twir:integration-token-refresh:" + providerName + ":" + channelID
}
