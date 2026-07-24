package vkvideo

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/platform"
	cfg "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/integrations/vk"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Config cfg.Config
}

type Provider struct {
	client *vk.OAuthClient
}

var _ platform.PlatformProvider = (*Provider)(nil)

func New(opts Opts) (*Provider, error) {
	client, err := vk.NewOAuthClient(vk.OAuthClientOpts{
		ClientID:      opts.Config.VKVideoClientID,
		ClientSecret:  opts.Config.VKVideoClientSecret,
		RedirectURL:   opts.Config.GetVkCallbackUrl(),
		APIBaseURL:    opts.Config.VKVideoAPIBaseURL,
		AuthBaseURL:   opts.Config.VKVideoAuthBaseURL,
		DevAPIBaseURL: opts.Config.VKVideoDevAPIBaseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("create VK Video OAuth client: %w", err)
	}

	return &Provider{client: client}, nil
}

func (p *Provider) Platform() platformentity.Platform {
	return platformentity.PlatformVKVideoLive
}

func (p *Provider) GetAuthURL(state, _ string) string {
	authorizationURL, err := p.client.AuthorizationURL(state, nil)
	if err != nil {
		return ""
	}

	return authorizationURL
}

func (p *Provider) ExchangeCode(
	ctx context.Context,
	input platform.ExchangeCodeInput,
) (*platform.PlatformTokens, error) {
	tokens, err := p.client.ExchangeCode(ctx, input.Code)
	if err != nil {
		return nil, fmt.Errorf("exchange VK Video code: %w", err)
	}

	return &platform.PlatformTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		Scopes:       tokens.Scopes,
	}, nil
}

func (p *Provider) RefreshToken(
	ctx context.Context,
	input platform.RefreshTokenInput,
) (*platform.PlatformTokens, error) {
	tokens, err := p.client.RefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh VK Video token: %w", err)
	}

	return &platform.PlatformTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		Scopes:       tokens.Scopes,
	}, nil
}

func (p *Provider) GetUser(ctx context.Context, accessToken string) (*platform.PlatformUser, error) {
	user, err := p.client.CurrentUser(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("get VK Video user: %w", err)
	}

	return &platform.PlatformUser{
		ID:          user.ID,
		Login:       user.Nick,
		DisplayName: user.Nick,
		Avatar:      user.Avatar,
	}, nil
}
