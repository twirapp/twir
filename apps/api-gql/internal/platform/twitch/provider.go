package twitchplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/kvizyx/twitchy/helix"
	twitchoauth "github.com/kvizyx/twitchy/oauth"
	"github.com/twirapp/twir/apps/api-gql/internal/platform"
	cfg "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"go.uber.org/fx"
)

const twitchOAuthBaseURL = "https://id.twitch.tv"

type Opts struct {
	fx.In

	Config cfg.Config
}

type Provider struct {
	config cfg.Config
}

var _ platform.PlatformProvider = (*Provider)(nil)

func New(opts Opts) *Provider {
	return &Provider{config: opts.Config}
}

func (p *Provider) Platform() platformentity.Platform {
	return platformentity.PlatformTwitch
}

func (p *Provider) newClient() (*twitchoauth.Client, error) {
	options := []twitchoauth.Option{twitchoauth.WithHTTPClient(&http.Client{})}
	if p.config.TwitchMockEnabled {
		options = append(options, twitchoauth.WithBaseURL(p.config.TwitchMockAuthUrl))
	}

	return twitchoauth.New(options...)
}

func (p *Provider) newHelixClient(accessToken string) (*helix.Client, error) {
	options := []helix.Option{
		helix.WithHTTPClient(&http.Client{}),
		helix.WithStaticToken(helix.Credential{
			AccessToken: accessToken,
			ClientID:    p.config.TwitchClientId,
			TokenClass:  helix.TokenClassUser,
		}),
	}
	if p.config.TwitchMockEnabled {
		options = append(options, helix.WithBaseURL(p.config.TwitchMockApiUrl))
	}

	return helix.New(options...)
}

func (p *Provider) GetAuthURL(state, _ string) string {
	authBaseURL := twitchOAuthBaseURL
	if p.config.TwitchMockEnabled {
		authBaseURL = p.config.TwitchMockAuthUrl
	}

	authorizeURL, err := url.Parse(authBaseURL)
	if err != nil {
		return ""
	}
	authorizeURL = authorizeURL.JoinPath("oauth2", "authorize")
	query := authorizeURL.Query()
	query.Set("response_type", "code")
	query.Set("client_id", p.config.TwitchClientId)
	query.Set("redirect_uri", p.config.GetTwitchCallbackUrl())
	query.Set("state", state)
	if p.config.TwitchMockEnabled {
		query.Set("scope", "")
	} else {
		query.Set("scope", "moderation:read channel:manage:broadcast channel:read:redemptions channel:manage:redemptions moderator:read:chatters moderator:manage:shoutouts moderator:manage:banned_users channel:read:vips channel:manage:vips channel:manage:moderators moderator:read:followers moderator:manage:chat_settings channel:read:polls channel:manage:polls channel:read:predictions channel:manage:predictions channel:read:subscriptions channel:moderate user:read:follows channel:bot channel:manage:raids")
	}
	authorizeURL.RawQuery = query.Encode()

	return authorizeURL.String()
}

func (p *Provider) ExchangeCode(ctx context.Context, input platform.ExchangeCodeInput) (*platform.PlatformTokens, error) {
	client, err := p.newClient()
	if err != nil {
		return nil, fmt.Errorf("create Twitch OAuth client: %w", err)
	}

	resp, err := client.ExchangeCode(ctx, twitchoauth.ExchangeCodeRequest{
		ClientID:     p.config.TwitchClientId,
		ClientSecret: p.config.TwitchClientSecret,
		Code:         input.Code,
		RedirectURI:  p.config.GetTwitchCallbackUrl(),
	})
	if err != nil {
		return nil, fmt.Errorf("exchange Twitch authorization code: %w", err)
	}

	return platformTokens(resp), nil
}

func (p *Provider) RefreshToken(ctx context.Context, input platform.RefreshTokenInput) (*platform.PlatformTokens, error) {
	client, err := p.newClient()
	if err != nil {
		return nil, fmt.Errorf("create Twitch OAuth client: %w", err)
	}

	resp, err := client.Refresh(ctx, twitchoauth.RefreshRequest{
		ClientID:     p.config.TwitchClientId,
		ClientSecret: p.config.TwitchClientSecret,
		RefreshToken: input.RefreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("refresh Twitch user access token: %w", err)
	}

	return platformTokens(resp), nil
}

func (p *Provider) GetUser(ctx context.Context, accessToken string) (*platform.PlatformUser, error) {
	client, err := p.newHelixClient(accessToken)
	if err != nil {
		return nil, fmt.Errorf("create helix client: %w", err)
	}

	resp, err := client.Users.GetUsers(ctx, helix.GetUsersRequest{})
	if err != nil {
		return nil, fmt.Errorf("get twitch users: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("twitch user not found")
	}

	u := resp.Data[0]

	return &platform.PlatformUser{
		ID:          u.ID,
		Login:       u.Login,
		DisplayName: u.DisplayName,
		Avatar:      u.ProfileImageURL,
	}, nil
}

func platformTokens(pair *twitchoauth.TokenPair) *platform.PlatformTokens {
	scopes := make([]string, len(pair.Scopes))
	for index, scope := range pair.Scopes {
		scopes[index] = string(scope)
	}

	return &platform.PlatformTokens{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    int(pair.ExpiresIn / time.Second),
		Scopes:       scopes,
	}
}
