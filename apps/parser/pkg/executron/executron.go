package executron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	config "github.com/twirapp/twir/libs/config"
)

type request struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type Response struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

func New(cfg config.Config, redisClient *redis.Client) Executron {
	return Executron{
		apiUrl:  cfg.ExecutronAddr,
		limiter: redis_rate.NewLimiter(redisClient),
	}
}

type Executron struct {
	apiUrl  string
	limiter *redis_rate.Limiter
}

func (c *Executron) ExecuteUserCode(
	ctx context.Context,
	channelId,
	language,
	code string,
) (*Response, error) {
	limitRes, err := c.limiter.Allow(ctx, "limits:executron:"+channelId, redis_rate.PerSecond(5))
	if err != nil {
		return nil, fmt.Errorf("cannot check rate limits for execute script: %w", err)
	}

	if limitRes.Allowed == 0 {
		return nil, fmt.Errorf("maximum 1 script execution per second per channel")
	}

	u, _ := url.Parse(c.apiUrl)
	u.Path = "/run"

	bodyData := request{
		Language: language,
		Code:     code,
	}
	bodyBytes, err := json.Marshal(bodyData)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cannot execute code: %s", string(respBody))
	}

	var executeResponse Response
	if err := json.Unmarshal(respBody, &executeResponse); err != nil {
		return nil, fmt.Errorf("cannot unmarshal response: %w", err)
	}

	if executeResponse.Error != "" {
		return nil, fmt.Errorf("cannot execute code: %s", executeResponse.Error)
	}

	return &executeResponse, nil
}
