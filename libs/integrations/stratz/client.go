package stratz

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL               = "https://api.stratz.com/graphql"
	maxResponseBodyBytes         = 1 << 20
	maxErrorResponsePreviewBytes = 4 << 10
)

var ErrDisabled = errors.New("stratz integration is disabled")

type NotablePlayer struct {
	AccountID int64
	Name      string
	TeamName  string
}

type Client struct {
	token      string
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

func New(token string, opts ...Option) *Client {
	c := &Client{
		token:   token,
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

func (c *Client) Enabled() bool {
	return c.token != ""
}

type graphQLError struct {
	Message string `json:"message"`
}

func (c *Client) do(ctx context.Context, query string, variables map[string]any, out any) error {
	if !c.Enabled() {
		return ErrDisabled
	}

	reqBody, err := json.Marshal(
		map[string]any{
			"query":     query,
			"variables": variables,
		},
	)
	if err != nil {
		return fmt.Errorf("stratz: failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return fmt.Errorf("stratz: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", "STRATZ_API")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("stratz: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyBytes+1))
	if err != nil {
		return fmt.Errorf("stratz: failed to read response body: %w", err)
	}
	if len(respBody) > maxResponseBodyBytes {
		return fmt.Errorf("stratz: response body exceeds %d bytes", maxResponseBodyBytes)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		preview := respBody
		if len(preview) > maxErrorResponsePreviewBytes {
			preview = preview[:maxErrorResponsePreviewBytes]
			return fmt.Errorf("stratz: unexpected status %d: %s...", resp.StatusCode, preview)
		}

		return fmt.Errorf("stratz: unexpected status %d: %s", resp.StatusCode, preview)
	}

	var envelope struct {
		Data   json.RawMessage `json:"data"`
		Errors []graphQLError  `json:"errors"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return fmt.Errorf("stratz: failed to parse response: %w", err)
	}

	if len(envelope.Errors) > 0 {
		return fmt.Errorf("stratz: graphql error: %s", envelope.Errors[0].Message)
	}

	if len(envelope.Data) == 0 {
		return fmt.Errorf("stratz: empty data in response")
	}

	if err := json.Unmarshal(envelope.Data, out); err != nil {
		return fmt.Errorf("stratz: failed to parse data: %w", err)
	}

	return nil
}

// Stratz live match schema is not publicly versioned; the query below targets
// live.match.liveWinRateValues which is the closest documented shape for live
// win rate values. Values can be probabilities or 0-100 percentages.
const winProbabilityQuery = `
query ($matchId: Long!) {
  live {
    match(id: $matchId) {
      liveWinRateValues { time winRate }
    }
  }
}`

func (c *Client) WinProbability(ctx context.Context, matchID int64) (float64, error) {
	var data struct {
		Live struct {
			Match *struct {
				LiveWinRateValues []struct {
					Time    int     `json:"time"`
					WinRate float64 `json:"winRate"`
				} `json:"liveWinRateValues"`
			} `json:"match"`
		} `json:"live"`
	}

	if err := c.do(ctx, winProbabilityQuery, map[string]any{"matchId": matchID}, &data); err != nil {
		return 0, err
	}

	if data.Live.Match == nil || len(data.Live.Match.LiveWinRateValues) == 0 {
		return 0, fmt.Errorf("stratz: no win rate values for match %d", matchID)
	}

	values := data.Live.Match.LiveWinRateValues
	winRate := values[len(values)-1].WinRate

	if winRate < 0 || winRate > 100 {
		return 0, fmt.Errorf("stratz: win rate %v is out of range for match %d", winRate, matchID)
	}

	for _, value := range values {
		if value.WinRate > 1 {
			return winRate / 100, nil
		}
	}

	return winRate, nil
}

// Stratz live match schema is not publicly versioned; players carry
// steamAccount.proSteamAccount when the account is a registered pro player,
// with the pro team reachable via proSteamAccount.team.tag.
const notablePlayersQuery = `
query ($matchId: Long!) {
  live {
    match(id: $matchId) {
      players {
        steamAccount {
          id
          name
          proSteamAccount {
            name
            team { tag }
          }
        }
      }
    }
  }
}`

func (c *Client) NotablePlayers(ctx context.Context, matchID int64) ([]NotablePlayer, error) {
	var data struct {
		Live struct {
			Match *struct {
				Players []struct {
					SteamAccount *struct {
						ID              int64  `json:"id"`
						Name            string `json:"name"`
						ProSteamAccount *struct {
							Name string `json:"name"`
							Team *struct {
								Tag string `json:"tag"`
							} `json:"team"`
						} `json:"proSteamAccount"`
					} `json:"steamAccount"`
				} `json:"players"`
			} `json:"match"`
		} `json:"live"`
	}

	if err := c.do(ctx, notablePlayersQuery, map[string]any{"matchId": matchID}, &data); err != nil {
		return nil, err
	}

	if data.Live.Match == nil {
		return nil, nil
	}

	notable := make([]NotablePlayer, 0, len(data.Live.Match.Players))
	for _, player := range data.Live.Match.Players {
		account := player.SteamAccount
		if account == nil || account.ProSteamAccount == nil {
			continue
		}

		np := NotablePlayer{
			AccountID: account.ID,
			Name:      account.ProSteamAccount.Name,
		}
		if np.Name == "" {
			np.Name = account.Name
		}
		if account.ProSteamAccount.Team != nil {
			np.TeamName = account.ProSteamAccount.Team.Tag
		}

		notable = append(notable, np)
	}

	return notable, nil
}
