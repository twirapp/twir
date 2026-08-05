package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Spotify) SearchTracks(ctx context.Context, query string, limit int) ([]SpotifyTrack, error) {
	if limit <= 0 {
		limit = 10
	}

	apiURL, err := url.Parse("https://api.spotify.com/v1/search")
	if err != nil {
		return nil, fmt.Errorf("failed to parse search url: %w", err)
	}
	q := apiURL.Query()
	q.Set("q", query)
	q.Set("type", "track")
	q.Set("limit", strconv.Itoa(limit))
	apiURL.RawQuery = q.Encode()

	resp, err := c.doRequest(ctx, http.MethodGet, apiURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, rateLimitedError(resp.Header)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cannot search spotify tracks: %s", string(resp.Body))
	}

	var data spotifySearchResponse
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	tracks := make([]SpotifyTrack, 0, len(data.Tracks.Items))
	for _, item := range data.Tracks.Items {
		tracks = append(tracks, trackFromItem(item))
	}

	return tracks, nil
}

func (c *Spotify) GetTrackByID(ctx context.Context, trackID string) (*SpotifyTrack, error) {
	if trackID == "" {
		return nil, ErrTrackNotFound
	}

	resp, err := c.doRequest(
		ctx,
		http.MethodGet,
		"https://api.spotify.com/v1/tracks/"+url.PathEscape(trackID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrTrackNotFound
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, rateLimitedError(resp.Header)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cannot get spotify track: %s", truncateBody(resp.Body))
	}

	var item spotifyTrackItem
	if err := json.Unmarshal(resp.Body, &item); err != nil {
		return nil, fmt.Errorf("failed to parse track response: %w", err)
	}

	track := trackFromItem(item)
	return &track, nil
}

func (c *Spotify) AddToQueue(ctx context.Context, trackURI string, deviceID string) error {
	apiURL, err := url.Parse("https://api.spotify.com/v1/me/player/queue")
	if err != nil {
		return fmt.Errorf("failed to parse queue url: %w", err)
	}
	q := apiURL.Query()
	q.Set("uri", trackURI)
	if deviceID != "" {
		q.Set("device_id", deviceID)
	}
	apiURL.RawQuery = q.Encode()

	resp, err := c.doRequest(ctx, http.MethodPost, apiURL.String(), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return playbackError(resp.StatusCode, resp.Body, resp.Header)
}

func (c *Spotify) GetCurrentlyPlaying(ctx context.Context) (*CurrentlyPlaying, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	if resp.StatusCode == http.StatusForbidden {
		return nil, ErrInsufficientScope
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, rateLimitedError(resp.Header)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cannot get currently playing track: %s", string(resp.Body))
	}

	var data spotifyCurrentlyPlayingResponse
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var item *SpotifyTrack
	if data.Item != nil {
		track := trackFromItem(*data.Item)
		item = &track
	}

	return &CurrentlyPlaying{
		CurrentlyPlayingType: data.CurrentlyPlayingType,
		Item:                 item,
		IsPlaying:            data.IsPlaying,
		ProgressMs:           data.ProgressMs,
		Device:               data.Device,
	}, nil
}

func (c *Spotify) GetQueue(ctx context.Context) ([]SpotifyTrack, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "https://api.spotify.com/v1/me/player/queue", nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, rateLimitedError(resp.Header)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cannot get spotify queue: %s", string(resp.Body))
	}

	var data spotifyQueueResponse
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	queue := make([]SpotifyTrack, 0, len(data.Queue))
	for _, item := range data.Queue {
		queue = append(queue, trackFromItem(item))
	}

	return queue, nil
}

func (c *Spotify) GetDevices(ctx context.Context) ([]Device, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "https://api.spotify.com/v1/me/player/devices", nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, rateLimitedError(resp.Header)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cannot get spotify devices: %s", string(resp.Body))
	}

	var data spotifyDevicesResponse
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return data.Devices, nil
}

func (c *Spotify) SkipNext(ctx context.Context, deviceID string) error {
	apiURL, err := url.Parse("https://api.spotify.com/v1/me/player/next")
	if err != nil {
		return fmt.Errorf("failed to parse next url: %w", err)
	}
	q := apiURL.Query()
	if deviceID != "" {
		q.Set("device_id", deviceID)
	}
	apiURL.RawQuery = q.Encode()

	resp, err := c.doRequest(ctx, http.MethodPost, apiURL.String(), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return playbackError(resp.StatusCode, resp.Body, resp.Header)
}
