package twitchplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/platform"
	cfg "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Config cfg.Config
}

type Provider struct {
	config cfg.Config
}

func New(opts Opts) *Provider {
	return &Provider{config: opts.Config}
}

func (p *Provider) Name() string {
	return platformentity.PlatformTwitch.String()
}

func (p *Provider) newClient() (*helix.Client, error) {
	httpClient := &http.Client{}
	apiBaseURL := ""
	if p.config.TwitchMockEnabled {
		apiBaseURL = p.config.TwitchMockApiUrl
	}

	return helix.NewClient(&helix.Options{
		ClientID:     p.config.TwitchClientId,
		ClientSecret: p.config.TwitchClientSecret,
		RedirectURI:  p.config.GetTwitchCallbackUrl(),
		HTTPClient:   httpClient,
		APIBaseURL:   apiBaseURL,
	})
}

func (p *Provider) GetAuthURL(state, _ string) string {
	client, err := p.newClient()
	if err != nil {
		return ""
	}

	return client.GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: "code",
		State:        state,
		Scopes: []string{
			"bits:read",
			"channel:bot",
			"channel:manage:broadcast",
			"channel:manage:moderators",
			"channel:manage:polls",
			"channel:manage:predictions",
			"channel:manage:raids",
			"channel:manage:redemptions",
			"channel:manage:schedule",
			"channel:manage:videos",
			"channel:manage:vips",
			"channel:moderate",
			"channel:read:goals",
			"channel:read:hype_train",
			"channel:read:polls",
			"channel:read:predictions",
			"channel:read:redemptions",
			"channel:read:subscriptions",
			"channel:read:vips",
			"chat:edit",
			"chat:read",
			"clips:edit",
			"moderation:read",
			"moderator:manage:announcements",
			"moderator:manage:automod",
			"moderator:manage:banned_users",
			"moderator:manage:blocked_terms",
			"moderator:manage:chat_messages",
			"moderator:manage:chat_settings",
			"moderator:manage:shield_mode",
			"moderator:manage:shoutouts",
			"moderator:manage:warnings",
			"moderator:read:chatters",
			"moderator:read:followers",
			"moderator:read:shield_mode",
			"moderator:read:shoutouts",
			"user:bot",
			"user:manage:blocked_users",
			"user:read:broadcast",
			"user:read:chat",
			"user:read:email",
			"user:read:follows",
			"user:read:subscriptions",
			"user:write:chat",
			"whispers:read",
		},
	})
}

func (p *Provider) ExchangeCode(ctx context.Context, code, _ string) (*platform.PlatformTokens, error) {
	client, err := p.newClient()
	if err != nil {
		return nil, fmt.Errorf("create helix client: %w", err)
	}

	resp, err := client.RequestUserAccessToken(code)
	if err != nil {
		return nil, fmt.Errorf("request user access token: %w", err)
	}

	if resp.ErrorMessage != "" {
		return nil, fmt.Errorf("twitch token exchange error: %s", resp.ErrorMessage)
	}

	return &platform.PlatformTokens{
		AccessToken:  resp.Data.AccessToken,
		RefreshToken: resp.Data.RefreshToken,
		ExpiresIn:    resp.Data.ExpiresIn,
		Scopes:       resp.Data.Scopes,
	}, nil
}

func (p *Provider) RefreshToken(ctx context.Context, refreshToken string) (*platform.PlatformTokens, error) {
	client, err := p.newClient()
	if err != nil {
		return nil, fmt.Errorf("create helix client: %w", err)
	}

	resp, err := client.RefreshUserAccessToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh user access token: %w", err)
	}

	if resp.ErrorMessage != "" {
		return nil, fmt.Errorf("twitch token refresh error: %s", resp.ErrorMessage)
	}

	return &platform.PlatformTokens{
		AccessToken:  resp.Data.AccessToken,
		RefreshToken: resp.Data.RefreshToken,
		ExpiresIn:    resp.Data.ExpiresIn,
		Scopes:       resp.Data.Scopes,
	}, nil
}

func (p *Provider) GetUser(ctx context.Context, accessToken string) (*platform.PlatformUser, error) {
	client, err := p.newClient()
	if err != nil {
		return nil, fmt.Errorf("create helix client: %w", err)
	}

	client.SetUserAccessToken(accessToken)

	resp, err := client.GetUsers(&helix.UsersParams{})
	if err != nil {
		return nil, fmt.Errorf("get twitch users: %w", err)
	}

	if resp.ErrorMessage != "" {
		return nil, fmt.Errorf("twitch get users error: %s", resp.ErrorMessage)
	}

	if len(resp.Data.Users) == 0 {
		return nil, fmt.Errorf("twitch user not found")
	}

	u := resp.Data.Users[0]

	return &platform.PlatformUser{
		ID:          u.ID,
		Login:       u.Login,
		DisplayName: u.DisplayName,
		Avatar:      u.ProfileImageURL,
	}, nil
}
