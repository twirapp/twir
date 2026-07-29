package nightbot

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

var ErrMissingClientCredentials = errors.New("nightbot OAuth client credentials are missing")

type ProviderError struct {
	StatusCode int
}

func (error ProviderError) Error() string {
	if error.StatusCode == 0 {
		return "Nightbot OAuth response is invalid"
	}
	return fmt.Sprintf("Nightbot OAuth request failed with status %d", error.StatusCode)
}

type refresher struct {
	client       *http.Client
	tokenURL     string
	integrations integrations.Repository
}

func (refresher refresher) Refresh(ctx context.Context, credential oauth.Credential) (oauth.RefreshResult, error) {
	integration, err := refresher.integrations.GetByService(ctx, integrationsmodel.ServiceNightbot)
	if err != nil {
		return oauth.RefreshResult{}, fmt.Errorf("load Nightbot integration settings: %w", err)
	}
	if integration.ClientID == nil || *integration.ClientID == "" || integration.ClientSecret == nil || *integration.ClientSecret == "" {
		return oauth.RefreshResult{}, ErrMissingClientCredentials
	}
	if credential.RefreshToken == "" {
		return oauth.RefreshResult{}, fmt.Errorf("Nightbot refresh token: %w", oauth.ErrInvalidCredential)
	}

	form := url.Values{
		"grant_type": {"refresh_token"}, "client_id": {*integration.ClientID},
		"client_secret": {*integration.ClientSecret}, "refresh_token": {credential.RefreshToken},
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, refresher.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return oauth.RefreshResult{}, fmt.Errorf("create Nightbot OAuth refresh request: %w", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := refresher.client.Do(request)
	if err != nil {
		return oauth.RefreshResult{}, fmt.Errorf("perform Nightbot OAuth refresh request: %w", err)
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
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil || payload.ExpiresIn <= 0 {
		return oauth.RefreshResult{}, ProviderError{}
	}
	result := oauth.RefreshResult{AccessToken: payload.AccessToken, ExpiresIn: time.Duration(payload.ExpiresIn) * time.Second}
	if payload.RefreshToken != "" {
		result.RefreshToken = &payload.RefreshToken
	}
	return result, nil
}
