package spotify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	"github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/model"
)

type Spotify struct {
	repo               channelsintegrationsspotify.Repository
	integration        deprecatedgormmodel.Integrations
	channelIntegration model.ChannelIntegrationSpotify
	isRetry            bool
}

func New(
	integration deprecatedgormmodel.Integrations,
	channelIntegration model.ChannelIntegrationSpotify,
	repo channelsintegrationsspotify.Repository,
) *Spotify {
	if channelIntegration.AccessToken == "" || channelIntegration.RefreshToken == "" {
		return nil
	}

	service := Spotify{
		integration:        integration,
		channelIntegration: channelIntegration,
		isRetry:            false,
		repo:               repo,
	}

	return &service
}

type spotifyRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (c *Spotify) refreshToken(ctx context.Context) error {
	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", c.channelIntegration.RefreshToken)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://accounts.spotify.com/api/token",
		bytes.NewBufferString(formData.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString(
		[]byte(c.integration.ClientID.String + ":" + c.integration.ClientSecret.String),
	)
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("cannot refresh spotify token: %s", string(body))
	}

	var data spotifyRefreshResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	input := channelsintegrationsspotify.UpdateInput{
		AccessToken: &data.AccessToken,
	}

	if data.RefreshToken != "" {
		input.RefreshToken = &data.RefreshToken
		c.channelIntegration.RefreshToken = data.RefreshToken
	}
	if err := c.repo.Update(ctx, c.channelIntegration.ID, input); err != nil {
		return fmt.Errorf("cannot save spotify token: %w", err)
	}

	c.channelIntegration.AccessToken = data.AccessToken

	return nil
}
