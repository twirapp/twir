package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/twirapp/twir/libs/oauth"
	integrations "github.com/twirapp/twir/libs/repositories/integrations"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

var ErrMissingClientCredentials = errors.New("spotify OAuth client credentials are missing")

type ProviderError struct {
	StatusCode int
}

func (error ProviderError) Error() string {
	if error.StatusCode == 0 {
		return "Spotify OAuth response is invalid"
	}
	return fmt.Sprintf("Spotify OAuth request failed with status %d", error.StatusCode)
}

type refresher struct {
	client       *http.Client
	tokenURL     string
	integrations integrations.Repository
}

func (refresher refresher) Refresh(ctx context.Context, credential oauth.Credential) (oauth.RefreshResult, error) {
	integration, err := refresher.integrations.GetByService(ctx, integrationsmodel.ServiceSpotify)
	if err != nil {
		return oauth.RefreshResult{}, fmt.Errorf("load Spotify integration settings: %w", err)
	}
	if integration.ClientID == nil || *integration.ClientID == "" || integration.ClientSecret == nil || *integration.ClientSecret == "" {
		return oauth.RefreshResult{}, ErrMissingClientCredentials
	}
	if credential.RefreshToken == "" {
		return oauth.RefreshResult{}, fmt.Errorf("Spotify refresh token: %w", oauth.ErrInvalidCredential)
	}

	form := url.Values{"grant_type": {"refresh_token"}, "refresh_token": {credential.RefreshToken}}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, refresher.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return oauth.RefreshResult{}, fmt.Errorf("create Spotify OAuth refresh request: %w", err)
	}
	request.SetBasicAuth(*integration.ClientID, *integration.ClientSecret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := refresher.client.Do(request)
	if err != nil {
		return oauth.RefreshResult{}, fmt.Errorf("perform Spotify OAuth refresh request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return oauth.RefreshResult{}, ProviderError{StatusCode: response.StatusCode}
	}

	var payload struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return oauth.RefreshResult{}, ProviderError{}
	}
	expiresIn := time.Duration(payload.ExpiresIn) * time.Second
	if expiresIn <= 0 {
		expiresIn = tokenLifetime
	}
	result := oauth.RefreshResult{AccessToken: payload.AccessToken, ExpiresIn: expiresIn}
	if payload.RefreshToken != "" {
		result.RefreshToken = &payload.RefreshToken
	}
	return result, nil
}
