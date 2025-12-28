package executron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	config "github.com/twirapp/twir/libs/config"
)

type request struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	UserId   string `json:"userId,omitempty"`
}

type Response struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

func New(cfg config.Config, redisClient *redis.Client) Executron {
	return Executron{
		isCf:           cfg.ExecutronCfClientId != "" && cfg.ExecutronCfClientSecret != "",
		apiUrl:         cfg.ExecutronAddr,
		cfClientId:     cfg.ExecutronCfClientId,
		cfClientSecret: cfg.ExecutronCfClientSecret,
		limiter:        redis_rate.NewLimiter(redisClient),
		redisClient:    redisClient,
	}
}

type Executron struct {
	isCf           bool
	apiUrl         string
	cfClientId     string
	cfClientSecret string

	limiter     *redis_rate.Limiter
	redisClient *redis.Client
}

func (c *Executron) executeUserCode(ctx context.Context, isCf bool, userId, language, code string) (*Response, error) {
	u, _ := url.Parse(c.apiUrl)
	if !isCf {
		u.Path = "/run"
	}

	bodyData := request{
		Language: language,
		Code:     code,
	}
	bodyBytes, err := json.Marshal(bodyData)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		u.String(),
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if isCf {
		storedCookie, err := c.redisClient.Get(ctx, "executron:cf:cookie").Result()
		if err == nil && storedCookie != "" {
			splitted := strings.Split(storedCookie, ";")
			splittedValue := strings.SplitN(splitted[0], "=", 2)
			req.Header.Set("cf-access-token", splittedValue[1])
		}

		req.Header.Set("CF-Access-Client-Id", c.cfClientId)
		req.Header.Set("CF-Access-Client-Secret", c.cfClientSecret)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if cfCookieHeader := resp.Header.Get("Set-Cookie"); isCf && cfCookieHeader != "" {
		if err := c.redisClient.Set(ctx, "executron:cf:cookie", cfCookieHeader, 0).Err(); err != nil {
			return nil, fmt.Errorf("cannot save cf cookie in cache: %w", err)
		}
	}

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

	return c.executeUserCode(ctx, c.isCf, channelId, language, code)
}
