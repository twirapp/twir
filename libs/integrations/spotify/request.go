package spotify

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type spotifyResponse struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

func (c *Spotify) doRequest(ctx context.Context, method, rawURL string, body io.Reader) (spotifyResponse, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return spotifyResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.channelIntegration.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return spotifyResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized && !c.isRetry && c.canRefresh() {
		c.isRetry = true
		defer func() {
			c.isRetry = false
		}()

		if err := c.refreshToken(ctx); err != nil {
			return spotifyResponse{}, err
		}

		return c.doRequest(ctx, method, rawURL, body)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return spotifyResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	return spotifyResponse{
		StatusCode: resp.StatusCode,
		Body:       bodyBytes,
		Header:     resp.Header.Clone(),
	}, nil
}
