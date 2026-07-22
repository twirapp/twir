package vkvideo

import (
	"context"
	"fmt"
	"strings"

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
	client *vk.IDClient
}

var _ platform.PlatformProvider = (*Provider)(nil)

func New(opts Opts) (*Provider, error) {
	client, err := vk.NewIDClient(vk.IDClientOpts{
		ClientID:     opts.Config.VKVideoClientID,
		ServiceToken: opts.Config.VKVideoServiceToken,
		RedirectURL:  opts.Config.VKVideoCallbackURL,
		APIBaseURL:   opts.Config.VKVideoAPIBaseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("create VK ID client: %w", err)
	}

	return &Provider{client: client}, nil
}

func (p *Provider) Name() string {
	return platformentity.PlatformVKVideoLive.String()
}

func (p *Provider) GetAuthURL(state, codeChallenge string) string {
	authorizationURL, err := p.client.AuthorizationURL(state, codeChallenge)
	if err != nil {
		return ""
	}

	return authorizationURL
}

func (p *Provider) ExchangeCode(
	ctx context.Context,
	input platform.ExchangeCodeInput,
) (*platform.PlatformTokens, error) {
	tokens, err := p.client.ExchangeCode(ctx, vk.IDExchangeCodeInput{
		Code:         input.Code,
		CodeVerifier: input.CodeVerifier,
		DeviceID:     input.DeviceID,
	})
	if err != nil {
		return nil, fmt.Errorf("exchange VK ID code: %w", err)
	}

	return &platform.PlatformTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		Scopes:       tokens.Scopes,
		DeviceID:     input.DeviceID,
	}, nil
}

func (p *Provider) RefreshToken(
	ctx context.Context,
	input platform.RefreshTokenInput,
) (*platform.PlatformTokens, error) {
	tokens, err := p.client.RefreshToken(ctx, vk.IDRefreshTokenInput{
		RefreshToken: input.RefreshToken,
		DeviceID:     input.DeviceID,
	})
	if err != nil {
		return nil, fmt.Errorf("refresh VK ID token: %w", err)
	}

	return &platform.PlatformTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		Scopes:       tokens.Scopes,
		DeviceID:     input.DeviceID,
	}, nil
}

func (p *Provider) GetUser(ctx context.Context, accessToken string) (*platform.PlatformUser, error) {
	user, err := p.client.UserInfo(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("get VK ID user: %w", err)
	}

	return &platform.PlatformUser{
		ID:          user.ID,
		DisplayName: strings.TrimSpace(user.FirstName + " " + user.LastName),
		Avatar:      user.Avatar,
	}, nil
}
