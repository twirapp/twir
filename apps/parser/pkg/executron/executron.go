package executron

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-redis/redis_rate/v10"
	"github.com/imroc/req/v3"
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
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

	var executeResponse Response
	resp, err := req.R().
		SetContext(ctx).
		SetBodyJsonMarshal(
			request{
				Language: language,
				Code:     code,
			},
		).
		SetSuccessResult(&executeResponse).
		Post(u.String())
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("cannot execute code: %s", resp.String())
	}

	if executeResponse.Error != "" {
		return nil, fmt.Errorf("cannot execute code: %s", executeResponse.Error)
	}

	return &executeResponse, nil
}
