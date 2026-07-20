package lastfm

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	lastfmAPIEndpoint  = "https://ws.audioscrobbler.com/2.0/"
	defaultHTTPTimeout = 10 * time.Second
	maxResponseSize    = int64(1 << 20)
	maxErrorBodySize   = int64(4 << 10)
)

type Opts struct {
	ApiKey       string
	ClientSecret string
	SessionKey   string
	UserName     string
	HTTPClient   *http.Client
}

type Lastfm struct {
	apiKey     string
	userName   string
	httpClient *http.Client
}

func New(ctx context.Context, opts Opts) (*Lastfm, error) {
	if opts.ApiKey == "" {
		return nil, errors.New("lastfm api key is required")
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHTTPTimeout}
	}

	lfm := &Lastfm{
		apiKey:     opts.ApiKey,
		userName:   opts.UserName,
		httpClient: httpClient,
	}
	if lfm.userName != "" {
		return lfm, nil
	}

	if opts.ClientSecret == "" || opts.SessionKey == "" {
		return nil, errors.New("lastfm client secret and session key are required when username is empty")
	}

	query := url.Values{}
	query.Set("method", "user.getinfo")
	query.Set("api_key", opts.ApiKey)
	query.Set("sk", opts.SessionKey)
	query.Set("api_sig", apiSignature(query, opts.ClientSecret))

	var response struct {
		User struct {
			Name string `json:"name"`
		} `json:"user"`
	}
	if err := lfm.requestJSON(ctx, http.MethodPost, query, &response); err != nil {
		return nil, fmt.Errorf("failed to resolve Last.fm username: %w", err)
	}
	if response.User.Name == "" {
		return nil, errors.New("lastfm user.getInfo returned an empty username")
	}

	lfm.userName = response.User.Name
	return lfm, nil
}

func apiSignature(params url.Values, secret string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		if key == "format" || key == "api_sig" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var value strings.Builder
	for _, key := range keys {
		value.WriteString(key)
		value.WriteString(params.Get(key))
	}
	value.WriteString(secret)

	signature := md5.Sum([]byte(value.String()))
	return fmt.Sprintf("%x", signature)
}

type apiErrorEnvelope struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

func (c *Lastfm) requestJSON(
	ctx context.Context,
	method string,
	params url.Values,
	target any,
) error {
	params.Set("api_key", c.apiKey)
	params.Set("format", "json")

	requestURL := lastfmAPIEndpoint
	var requestBody io.Reader
	if method == http.MethodPost {
		requestBody = strings.NewReader(params.Encode())
	} else {
		requestURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		requestURL,
		requestBody,
	)
	if err != nil {
		return fmt.Errorf("failed to create Last.fm request: %w", err)
	}
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform Last.fm request: %w", err)
	}
	defer resp.Body.Close()

	isSuccess := resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
	responseLimit := maxResponseSize
	if !isSuccess {
		responseLimit = maxErrorBodySize
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, responseLimit+1))
	if err != nil {
		return fmt.Errorf("failed to read Last.fm response: %w", err)
	}
	overSize := int64(len(body)) > responseLimit
	if overSize {
		body = body[:responseLimit]
	}

	var apiError apiErrorEnvelope
	if err := json.Unmarshal(body, &apiError); err == nil && apiError.Error != 0 {
		message := boundedErrorMessage(strings.TrimSpace(apiError.Message))
		if !isSuccess {
			return fmt.Errorf(
				"Last.fm API returned %s with error %d: %s",
				resp.Status,
				apiError.Error,
				message,
			)
		}
		return fmt.Errorf("Last.fm API error %d: %s", apiError.Error, message)
	}

	if !isSuccess {
		message := strings.TrimSpace(string(body))
		if message == "" {
			message = "empty response body"
		}
		if overSize {
			message += "... (truncated)"
		}
		return fmt.Errorf("Last.fm API returned %s: %s", resp.Status, message)
	}
	if overSize {
		return fmt.Errorf("Last.fm response exceeds %d bytes", maxResponseSize)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to decode Last.fm response: %w", err)
	}

	return nil
}

func boundedErrorMessage(message string) string {
	if message == "" {
		return "unknown error"
	}
	if int64(len(message)) <= maxErrorBodySize {
		return message
	}
	return message[:maxErrorBodySize] + "... (truncated)"
}

type Track struct {
	Title     string
	Artist    string
	Image     string
	PlayedUTS string
}

type recentTracksResponse struct {
	RecentTracks struct {
		Tracks []struct {
			Name   string `json:"name"`
			Artist struct {
				Name string `json:"#text"`
			} `json:"artist"`
			Images []struct {
				URL string `json:"#text"`
			} `json:"image"`
			Attributes struct {
				NowPlaying string `json:"nowplaying"`
			} `json:"@attr"`
			Date struct {
				UTS  string `json:"uts"`
				Text string `json:"#text"`
			} `json:"date"`
		} `json:"track"`
	} `json:"recenttracks"`
}

func (c *Lastfm) getRecentTracks(ctx context.Context, limit int) (*recentTracksResponse, error) {
	query := url.Values{}
	query.Set("method", "user.getrecenttracks")
	query.Set("user", c.userName)
	query.Set("limit", strconv.Itoa(limit))

	var response recentTracksResponse
	if err := c.requestJSON(ctx, http.MethodGet, query, &response); err != nil {
		return nil, fmt.Errorf("failed to get Last.fm recent tracks: %w", err)
	}

	return &response, nil
}

func (c *Lastfm) GetTrack(ctx context.Context) (*Track, error) {
	response, err := c.getRecentTracks(ctx, 1)
	if err != nil {
		return nil, err
	}
	if len(response.RecentTracks.Tracks) == 0 ||
		response.RecentTracks.Tracks[0].Attributes.NowPlaying != "true" {
		return nil, nil
	}

	track := response.RecentTracks.Tracks[0]
	var cover string
	if len(track.Images) > 0 {
		cover = track.Images[0].URL
	}

	return &Track{
		Title:     track.Name,
		Artist:    track.Artist.Name,
		Image:     cover,
		PlayedUTS: track.Date.Text,
	}, nil
}

func (c *Lastfm) GetRecentTracks(ctx context.Context, limit int) ([]Track, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	response, err := c.getRecentTracks(ctx, limit)
	if err != nil {
		return nil, err
	}

	recentTracks := make([]Track, 0, len(response.RecentTracks.Tracks))
	for _, track := range response.RecentTracks.Tracks {
		if track.Attributes.NowPlaying == "true" {
			continue
		}

		var cover string
		if len(track.Images) > 0 {
			cover = track.Images[0].URL
		}

		recentTracks = append(recentTracks, Track{
			Title:     track.Name,
			Artist:    track.Artist.Name,
			Image:     cover,
			PlayedUTS: track.Date.UTS,
		})
	}

	return recentTracks, nil
}
