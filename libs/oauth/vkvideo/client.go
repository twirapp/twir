package vkvideo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/twirapp/twir/libs/integrations/vk"
	"github.com/twirapp/twir/libs/oauth"
)

const Provider = oauth.Provider("vk_video")

var ErrMissingRefreshToken = errors.New("VK Video refresh token is missing")

type Client interface {
	RefreshToken(context.Context, string) (*vk.OAuthToken, error)
}

type ClientFactory func() (Client, error)

type ProviderError struct {
	StatusCode int
}

func (error ProviderError) Error() string {
	if error.StatusCode > 0 {
		return fmt.Sprintf("VK Video OAuth request failed with status %d", error.StatusCode)
	}
	return "VK Video OAuth request failed"
}

type refresher struct {
	factory ClientFactory
}

func (refresher refresher) Refresh(ctx context.Context, credential oauth.Credential) (oauth.RefreshResult, error) {
	if credential.RefreshToken == "" {
		return oauth.RefreshResult{}, fmt.Errorf("%w: %w", ErrMissingRefreshToken, oauth.ErrInvalidCredential)
	}

	client, err := refresher.factory()
	if err != nil {
		return oauth.RefreshResult{}, fmt.Errorf("create VK Video OAuth client: %w", err)
	}
	response, err := client.RefreshToken(ctx, credential.RefreshToken)
	if err != nil {
		return oauth.RefreshResult{}, redactProviderError(err)
	}

	result := oauth.RefreshResult{
		AccessToken: response.AccessToken,
		ExpiresIn:   time.Duration(response.ExpiresIn) * time.Second,
	}
	if response.RefreshToken != "" {
		refreshToken := response.RefreshToken
		result.RefreshToken = &refreshToken
	}
	if len(response.Scopes) > 0 {
		result.Scopes = append([]string(nil), response.Scopes...)
	}
	return result, nil
}

func redactProviderError(err error) error {
	var providerError *vk.ProviderError
	if errors.As(err, &providerError) {
		return ProviderError{StatusCode: providerError.StatusCode}
	}
	return fmt.Errorf("VK Video OAuth request: %w", err)
}
