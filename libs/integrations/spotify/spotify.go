package spotify

import (
	"context"
	"fmt"

	deprecatedgormmodel "github.com/satont/twir/libs/gomodels"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	"github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/model"

	"github.com/imroc/req/v3"
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
	data := spotifyRefreshResponse{}

	body := make(map[string]string, 2)
	body["grant_type"] = "refresh_token"
	body["refresh_token"] = c.channelIntegration.RefreshToken

	resp, err := req.R().
		SetContext(ctx).
		SetFormData(body).
		SetSuccessResult(&data).
		SetBasicAuth(
			c.integration.ClientID.String,
			c.integration.ClientSecret.String,
		).
		Post("https://accounts.spotify.com/api/token")

	if err != nil {
		return err
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("cannot refresh spotify token: %s", resp.String())
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
