package platform

import (
	"context"

	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type PlatformTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	Scopes       []string
	DeviceID     string
}

type ExchangeCodeInput struct {
	Code         string
	CodeVerifier string
	DeviceID     string
}

type RefreshTokenInput struct {
	RefreshToken string
	DeviceID     string
}

type PlatformUser struct {
	ID          string
	Login       string
	DisplayName string
	Avatar      string
}

type PlatformProvider interface {
	Platform() platformentity.Platform
	GetAuthURL(state, codeChallenge string) string
	ExchangeCode(ctx context.Context, input ExchangeCodeInput) (*PlatformTokens, error)
	RefreshToken(ctx context.Context, input RefreshTokenInput) (*PlatformTokens, error)
	GetUser(ctx context.Context, accessToken string) (*PlatformUser, error)
}
