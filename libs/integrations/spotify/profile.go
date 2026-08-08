package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	resp, err := c.doRequest(ctx, http.MethodGet, "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cannot get profile: %s", string(resp.Body))
	}

	var data SpotifyProfile
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &data, nil
}
