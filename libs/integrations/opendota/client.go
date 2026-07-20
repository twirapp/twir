package opendota

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultBaseURL = "https://api.opendota.com/api"

type RecentMatch struct {
	MatchID    int64 `json:"match_id"`
	HeroID     int   `json:"hero_id"`
	Kills      int   `json:"kills"`
	Deaths     int   `json:"deaths"`
	Assists    int   `json:"assists"`
	Duration   int   `json:"duration"`
	PlayerSlot int   `json:"player_slot"`
	RadiantWin bool  `json:"radiant_win"`
	GameMode   int   `json:"game_mode"`
	LobbyType  int   `json:"lobby_type"`
	StartTime  int64 `json:"start_time"`
}

func (m RecentMatch) IsRadiant() bool {
	return m.PlayerSlot < 128
}

func (m RecentMatch) Won() bool {
	return m.IsRadiant() == m.RadiantWin
}

type ProPlayer struct {
	AccountID int64  `json:"account_id"`
	Name      string `json:"name"`
	TeamName  string `json:"team_name"`
	TeamTag   string `json:"team_tag"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func New(opts ...Option) *Client {
	c := &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) get(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("opendota: failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("opendota: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("opendota: failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opendota: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("opendota: failed to parse response: %w", err)
	}

	return nil
}

func (c *Client) RecentMatches(ctx context.Context, accountID int64) ([]RecentMatch, error) {
	var matches []RecentMatch
	if err := c.get(ctx, fmt.Sprintf("/players/%d/recentMatches", accountID), &matches); err != nil {
		return nil, err
	}

	return matches, nil
}

func (c *Client) Heroes(ctx context.Context) (map[int]string, error) {
	var heroes map[int]struct {
		LocalizedName string `json:"localized_name"`
	}
	if err := c.get(ctx, "/constants/heroes", &heroes); err != nil {
		return nil, err
	}

	result := make(map[int]string, len(heroes))
	for id, hero := range heroes {
		result[id] = hero.LocalizedName
	}

	return result, nil
}

func (c *Client) ProPlayers(ctx context.Context) ([]ProPlayer, error) {
	var players []ProPlayer
	if err := c.get(ctx, "/proPlayers", &players); err != nil {
		return nil, err
	}

	return players, nil
}
