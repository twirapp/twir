package nightbot

import (
	"context"
	"fmt"
	"time"

	"github.com/twirapp/twir/libs/oauth"
	channelsintegrations "github.com/twirapp/twir/libs/repositories/channels_integrations"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

type store struct {
	repository channelsintegrations.Repository
}

func (store store) Load(ctx context.Context, key oauth.CredentialKey) (oauth.Credential, error) {
	integration, err := store.repository.GetByChannelAndService(ctx, string(key.ID), integrationsmodel.ServiceNightbot)
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("load Nightbot channel integration: %w", err)
	}
	accessToken := ""
	if integration.AccessToken != nil {
		accessToken = *integration.AccessToken
	}
	refreshToken := ""
	if integration.RefreshToken != nil {
		refreshToken = *integration.RefreshToken
	}
	return oauth.Credential{
		Provider: key.Provider, ID: key.ID, AccessToken: accessToken, RefreshToken: refreshToken,
		ObtainedAt: time.Unix(0, 0), ExpiresIn: time.Nanosecond,
	}, nil
}

func (store store) Commit(ctx context.Context, credential oauth.Credential) error {
	integration, err := store.repository.GetByChannelAndService(ctx, string(credential.ID), integrationsmodel.ServiceNightbot)
	if err != nil {
		return fmt.Errorf("load Nightbot channel integration for commit: %w", err)
	}
	enabled := true
	input := channelsintegrations.UpdateInput{Enabled: &enabled, AccessToken: &credential.AccessToken}
	if integration.RefreshToken == nil || *integration.RefreshToken != credential.RefreshToken {
		input.RefreshToken = &credential.RefreshToken
	}
	if err := store.repository.Update(ctx, integration.ID, input); err != nil {
		return fmt.Errorf("persist Nightbot channel integration: %w", err)
	}
	return nil
}
