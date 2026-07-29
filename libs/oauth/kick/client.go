package kick

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/scorfly/gokick"
	"github.com/twirapp/twir/libs/oauth"
)

const Provider = oauth.Provider("kick")

type Client interface {
	RefreshToken(context.Context, string) (gokick.TokenResponse, error)
}

type ClientFactory func() (Client, error)

type AppClient interface {
	GetAppAccessToken(context.Context) (gokick.AppTokenResponse, error)
}

type AppClientFactory func() (AppClient, error)

type ProviderError struct {
	StatusCode int
}

func (error ProviderError) Error() string {
	return fmt.Sprintf("kick OAuth request failed with status %d", error.StatusCode)
}

func newRefresher(factory ClientFactory) oauth.Refresher {
	return refresher{factory: factory}
}

type refresher struct {
	factory ClientFactory
}

func (refresher refresher) Refresh(ctx context.Context, credential oauth.Credential) (oauth.RefreshResult, error) {
	client, err := refresher.factory()
	if err != nil {
		return oauth.RefreshResult{}, fmt.Errorf("create Kick OAuth client: %w", err)
	}

	response, err := client.RefreshToken(ctx, credential.RefreshToken)
	if err != nil {
		return oauth.RefreshResult{}, redactProviderError(err)
	}

	result := oauth.RefreshResult{
		AccessToken: response.AccessToken,
		Scopes:      strings.Fields(response.Scope),
		ExpiresIn:   time.Duration(response.ExpiresIn) * time.Second,
	}
	if response.RefreshToken != "" {
		refreshToken := response.RefreshToken
		result.RefreshToken = &refreshToken
	}

	return result, nil
}

type appFetcher struct {
	factory AppClientFactory
}

func (fetcher appFetcher) FetchAppToken(ctx context.Context, _ oauth.AppTokenKey) (oauth.AppToken, error) {
	client, err := fetcher.factory()
	if err != nil {
		return oauth.AppToken{}, fmt.Errorf("create Kick OAuth client: %w", err)
	}

	response, err := client.GetAppAccessToken(ctx)
	if err != nil {
		return oauth.AppToken{}, redactProviderError(err)
	}

	return oauth.AppToken{
		AccessToken: response.AccessToken,
		ObtainedAt:  time.Now().UTC(),
		ExpiresIn:   time.Duration(response.ExpiresIn) * time.Second,
	}, nil
}

func redactProviderError(err error) error {
	var providerError interface {
		error
		Code() int
	}
	if errors.As(err, &providerError) && providerError.Code() >= 400 && providerError.Code() < 500 {
		return ProviderError{StatusCode: providerError.Code()}
	}
	return fmt.Errorf("Kick OAuth request: %w", err)
}
