package spotify

import (
	"context"
	"fmt"
	"time"

	"github.com/twirapp/twir/libs/oauth"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
)

const tokenLifetime = time.Hour

type store struct {
	repository channelsintegrationsspotify.Repository
}

func (store store) Load(ctx context.Context, key oauth.CredentialKey) (oauth.Credential, error) {
	integration, err := store.repository.GetByChannelID(ctx, string(key.ID))
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("load Spotify channel integration: %w", err)
	}
	obtainedAt := integration.UpdatedAt
	if obtainedAt.IsZero() {
		obtainedAt = time.Unix(0, 0)
	}
	return oauth.Credential{
		Provider: key.Provider, ID: key.ID, AccessToken: integration.AccessToken, RefreshToken: integration.RefreshToken,
		Scopes: integration.Scopes, ObtainedAt: obtainedAt, ExpiresIn: tokenLifetime,
	}, nil
}

func (store store) Commit(ctx context.Context, credential oauth.Credential) error {
	integration, err := store.repository.GetByChannelID(ctx, string(credential.ID))
	if err != nil {
		return fmt.Errorf("load Spotify channel integration for commit: %w", err)
	}
	input := channelsintegrationsspotify.UpdateInput{AccessToken: &credential.AccessToken}
	if integration.RefreshToken != credential.RefreshToken {
		input.RefreshToken = &credential.RefreshToken
	}
	if err := store.repository.Update(ctx, integration.ID, input); err != nil {
		return fmt.Errorf("persist Spotify channel integration: %w", err)
	}
	return nil
}
