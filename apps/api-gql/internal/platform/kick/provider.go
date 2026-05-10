package kickplatform

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/scorfly/gokick"
	"github.com/twirapp/twir/apps/api-gql/internal/platform"
	cfg "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"go.uber.org/fx"
)

var scopes = []gokick.Scope{
	gokick.ScopeUserRead,
	gokick.ScopeEventSubscribe,
	gokick.ScopeChatWrite,
	gokick.ScopeChannelRead,
}

type Opts struct {
	fx.In

	Config cfg.Config
}

type Provider struct {
	config cfg.Config
	client *gokick.Client
}

func New(opts Opts) *Provider {
	client, _ := gokick.NewClient(
		&gokick.ClientOptions{
			ClientID:     opts.Config.KickClientId,
			ClientSecret: opts.Config.KickClientSecret,
		})

	return &Provider{
		config: opts.Config,
		client: client,
	}
}

func (p *Provider) Name() string {
	return platformentity.PlatformKick.String()
}

func (p *Provider) GetAuthURL(state, codeChallenge string) string {
	return p.buildAuthURL(state, codeChallenge, p.config.GetKickCallbackUrl())
}

func (p *Provider) GetBotSetupAuthURL(state, codeChallenge string) string {
	u, err := url.Parse(p.config.SiteBaseUrl)
	if err != nil {
		panic(err)
	}

	return p.buildAuthURL(state, codeChallenge, u.JoinPath("api", "auth", "kick", "bot-callback").String())
}

func (p *Provider) buildAuthURL(state, codeChallenge, redirectURI string) string {
	authURL, _ := p.client.GetAuthorize(
		redirectURI,
		state,
		codeChallenge,
		scopes,
	)
	return authURL
}

func (p *Provider) ExchangeCode(ctx context.Context, code, codeVerifier string) (*platform.PlatformTokens, error) {
	return p.exchangeCodeWithRedirectURI(ctx, code, codeVerifier, p.config.GetKickCallbackUrl())
}

func (p *Provider) ExchangeBotSetupCode(ctx context.Context, code, codeVerifier string) (*platform.PlatformTokens, error) {
	u, err := url.Parse(p.config.SiteBaseUrl)
	if err != nil {
		return nil, err
	}
	redirectURI := u.JoinPath("api", "auth", "kick", "bot-callback").String()
	return p.exchangeCodeWithRedirectURI(ctx, code, codeVerifier, redirectURI)
}

func (p *Provider) exchangeCodeWithRedirectURI(ctx context.Context, code, codeVerifier, redirectURI string) (*platform.PlatformTokens, error) {
	token, err := p.client.GetToken(ctx, redirectURI, code, codeVerifier)
	if err != nil {
		return nil, fmt.Errorf("kick token exchange: %w", err)
	}

	return &platform.PlatformTokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
		Scopes:       parseScopes(token.Scope),
	}, nil
}

func (p *Provider) RefreshToken(ctx context.Context, refreshToken string) (*platform.PlatformTokens, error) {
	token, err := p.client.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("kick token refresh: %w", err)
	}

	return &platform.PlatformTokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
		Scopes:       parseScopes(token.Scope),
	}, nil
}

func (p *Provider) GetUser(ctx context.Context, accessToken string) (*platform.PlatformUser, error) {
	client, err := gokick.NewClient(
		&gokick.ClientOptions{
			UserAccessToken: accessToken,
		})
	if err != nil {
		return nil, fmt.Errorf("create kick client: %w", err)
	}

	response, err := client.GetUsers(ctx, gokick.NewUserListFilter())
	if err != nil {
		return nil, fmt.Errorf("kick get user: %w", err)
	}

	if len(response.Result) == 0 {
		return nil, fmt.Errorf("kick user not found")
	}

	user := response.Result[0]

	return &platform.PlatformUser{
		ID:          strconv.Itoa(user.UserID),
		Login:       user.Name,
		DisplayName: user.Name,
		Avatar:      user.ProfilePicture,
	}, nil
}

func parseScopes(scope string) []string {
	if scope == "" {
		return nil
	}
	return strings.Fields(scope)
}
