package platform

import "context"

type PlatformTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	Scopes       []string
}

type PlatformUser struct {
	ID          string
	Login       string
	DisplayName string
	Avatar      string
}

type PlatformProvider interface {
	Name() string
	GetAuthURL(state, codeChallenge string) string
	ExchangeCode(ctx context.Context, code, codeVerifier string) (*PlatformTokens, error)
	RefreshToken(ctx context.Context, refreshToken string) (*PlatformTokens, error)
	GetUser(ctx context.Context, accessToken string) (*PlatformUser, error)
}
