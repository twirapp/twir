package spotify

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
)

type SpotifyProfile struct {
	Country      string `json:"country"`
	DisplayName  string `json:"display_name"`
	Email        string `json:"email"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href      string `json:"href"`
	ID        string `json:"id"`
	Product   string `json:"product"`
	Type      string `json:"type"`
	URI       string `json:"uri"`
	Followers struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`
}

func (c *Spotify) GetProfile(ctx context.Context) (*SpotifyProfile, error) {
	data := SpotifyProfile{}
	resp, err := req.R().
		SetContext(ctx).
		SetBearerAuthToken(c.channelIntegration.AccessToken).
		SetSuccessResult(&data).
		Get("https://api.spotify.com/v1/me")

	if resp.StatusCode == 401 && !c.isRetry {
		c.isRetry = true
		c.refreshToken(ctx)
		return c.GetProfile(ctx)
	}

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("cannot get profile: %s", resp.String())
	}

	return &data, nil
}
