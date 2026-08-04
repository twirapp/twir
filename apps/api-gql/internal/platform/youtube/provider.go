package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/twirapp/twir/apps/api-gql/internal/platform"
	cfg "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"golang.org/x/oauth2"
)

const (
	youtubeAuthorizationURL = "https://accounts.google.com/o/oauth2/v2/auth"
	youtubeTokenURL         = "https://oauth2.googleapis.com/token"
	youtubeChannelsURL      = "https://www.googleapis.com/youtube/v3/channels?part=snippet&mine=true"
	youtubeScope            = "https://www.googleapis.com/auth/youtube"
)

type Opts struct {
	Config cfg.Config
}

type Provider struct {
	config     cfg.Config
	httpClient *http.Client
}

type youtubeThumbnail struct {
	URL string `json:"url"`
}

var _ platform.PlatformProvider = (*Provider)(nil)

func New(opts Opts) *Provider {
	return &Provider{config: opts.Config, httpClient: http.DefaultClient}
}

func (p *Provider) Platform() platformentity.Platform {
	return platformentity.PlatformYouTube
}

func (p *Provider) GetAuthURL(state, codeChallenge string) string {
	return p.authorizationURL(state, codeChallenge, p.config.GetYouTubeCallbackUrl())
}

func (p *Provider) GetBotSetupAuthURL(state, codeChallenge string) string {
	return p.authorizationURL(state, codeChallenge, p.config.GetYouTubeBotCallbackUrl())
}

func (p *Provider) ExchangeCode(ctx context.Context, input platform.ExchangeCodeInput) (*platform.PlatformTokens, error) {
	return p.exchangeCode(ctx, input.Code, input.CodeVerifier, p.config.GetYouTubeCallbackUrl())
}

func (p *Provider) ExchangeBotSetupCode(ctx context.Context, code, codeVerifier string) (*platform.PlatformTokens, error) {
	return p.exchangeCode(ctx, code, codeVerifier, p.config.GetYouTubeBotCallbackUrl())
}

func (p *Provider) RefreshToken(ctx context.Context, input platform.RefreshTokenInput) (*platform.PlatformTokens, error) {
	oauthConfig := p.oauthConfig(p.config.GetYouTubeCallbackUrl())
	token, err := oauthConfig.TokenSource(ctx, &oauth2.Token{RefreshToken: input.RefreshToken}).Token()
	if err != nil {
		return nil, fmt.Errorf("refresh YouTube token: %w", err)
	}
	if token.RefreshToken == "" {
		token.RefreshToken = input.RefreshToken
	}

	return platformTokensFromOAuthToken(token), nil
}

func (p *Provider) GetUser(ctx context.Context, accessToken string) (*platform.PlatformUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, youtubeChannelsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create YouTube channels request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get YouTube channel: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get YouTube channel: unexpected status %s", resp.Status)
	}

	var result struct {
		Items []struct {
			ID      string `json:"id"`
			Snippet struct {
				Title      string                      `json:"title"`
				CustomURL  string                      `json:"customUrl"`
				Thumbnails map[string]youtubeThumbnail `json:"thumbnails"`
			} `json:"snippet"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode YouTube channel: %w", err)
	}
	if len(result.Items) == 0 {
		return nil, fmt.Errorf("YouTube channel not found")
	}

	channel := result.Items[0]
	return &platform.PlatformUser{ID: channel.ID, Login: channel.Snippet.CustomURL, DisplayName: channel.Snippet.Title, Avatar: youtubeAvatar(channel.Snippet.Thumbnails)}, nil
}

func (p *Provider) authorizationURL(state, codeChallenge, redirectURL string) string {
	oauthConfig := p.oauthConfig(redirectURL)
	return oauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)
}

func (p *Provider) exchangeCode(ctx context.Context, code, codeVerifier, redirectURL string) (*platform.PlatformTokens, error) {
	oauthConfig := p.oauthConfig(redirectURL)
	token, err := oauthConfig.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier))
	if err != nil {
		return nil, fmt.Errorf("exchange YouTube code: %w", err)
	}

	return platformTokensFromOAuthToken(token), nil
}

func (p *Provider) oauthConfig(redirectURL string) oauth2.Config {
	return oauth2.Config{ClientID: p.config.YouTubeClientID, ClientSecret: p.config.YouTubeClientSecret, RedirectURL: redirectURL, Scopes: []string{youtubeScope}, Endpoint: oauth2.Endpoint{AuthURL: youtubeAuthorizationURL, TokenURL: youtubeTokenURL}}
}

func platformTokensFromOAuthToken(token *oauth2.Token) *platform.PlatformTokens {
	expiresIn := int(time.Until(token.Expiry).Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}
	scope, _ := token.Extra("scope").(string)
	return &platform.PlatformTokens{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken, ExpiresIn: expiresIn, Scopes: strings.Fields(scope)}
}

func youtubeAvatar(thumbnails map[string]youtubeThumbnail) string {
	for _, size := range []string{"maxres", "standard", "high", "medium", "default"} {
		if thumbnail := thumbnails[size]; thumbnail.URL != "" {
			return thumbnail.URL
		}
	}

	return ""
}
