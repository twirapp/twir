package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.channelIntegration.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 && !c.isRetry {
		c.isRetry = true
		c.refreshToken(ctx)
		return c.GetProfile(ctx)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cannot get profile: %s", string(body))
	}

	var data SpotifyProfile
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &data, nil
}
