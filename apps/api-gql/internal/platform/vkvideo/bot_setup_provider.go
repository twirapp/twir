package vkvideo

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/platform"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/vk"
	"go.uber.org/fx"
)

type BotSetupProviderOpts struct {
	fx.In

	Config cfg.Config
}

type BotSetupProvider struct {
	config cfg.Config
}

// POST /v1/chat/message/send requires user authorization with this permission.
const vkVideoBotChatSendScope = "chat:message:send"

func NewBotSetupProvider(opts BotSetupProviderOpts) *BotSetupProvider {
	return &BotSetupProvider{config: opts.Config}
}

func (p *BotSetupProvider) GetBotSetupAuthURL(state string) (string, error) {
	client, err := p.newClient()
	if err != nil {
		return "", err
	}

	return client.AuthorizationURL(state, []string{vkVideoBotChatSendScope})
}

func (p *BotSetupProvider) ExchangeBotSetupCode(
	ctx context.Context,
	code string,
) (*platform.PlatformTokens, error) {
	client, err := p.newClient()
	if err != nil {
		return nil, err
	}

	tokens, err := client.ExchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange VK Video bot code: %w", err)
	}

	return &platform.PlatformTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		Scopes:       tokens.Scopes,
	}, nil
}

func (p *BotSetupProvider) GetUser(ctx context.Context, accessToken string) (*platform.PlatformUser, error) {
	client, err := p.newClient()
	if err != nil {
		return nil, err
	}

	user, err := client.CurrentUser(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("get VK Video bot user: %w", err)
	}

	return &platform.PlatformUser{ID: user.ID, Login: user.Nick, DisplayName: user.Nick, Avatar: user.Avatar}, nil
}

func (p *BotSetupProvider) newClient() (*vk.OAuthClient, error) {
	client, err := vk.NewOAuthClient(vk.OAuthClientOpts{
		ClientID:      p.config.VKVideoClientID,
		ClientSecret:  p.config.VKVideoClientSecret,
		RedirectURL:   p.config.GetVkVideoBotCallbackUrl(),
		APIBaseURL:    p.config.VKVideoAPIBaseURL,
		AuthBaseURL:   p.config.VKVideoAuthBaseURL,
		DevAPIBaseURL: p.config.VKVideoDevAPIBaseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("create VK Video bot OAuth client: %w", err)
	}

	return client, nil
}
