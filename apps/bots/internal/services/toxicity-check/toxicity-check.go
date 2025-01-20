package toxicity_check

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/imroc/req/v3"
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

type Opts struct {
	fx.In

	RedisClient *redis.Client
	Config      config.Config
}

func New(opts Opts) *Service {
	return &Service{
		redisClient: opts.RedisClient,
		config:      opts.Config,
	}
}

type Service struct {
	redisClient *redis.Client
	config      config.Config
}

func (c *Service) CheckTextToxicity(ctx context.Context, text string) (bool, error) {
	if c.config.ToxicityAddr == "" {
		if c.config.AppEnv == "production" {
			return false, fmt.Errorf("toxicity addr is not set")
		}

		return false, nil
	}

	textHash := sha256.Sum256([]byte(text))
	redisKey := fmt.Sprintf("toxicity-check:%s", fmt.Sprintf("%x", textHash))

	cachedToxicity, err := c.redisClient.Get(ctx, redisKey).Bool()
	cachedExists := !errors.Is(err, redis.Nil)
	if err != nil && cachedExists {
		return false, fmt.Errorf("cannot get cached toxicity: %w", err)
	}

	if cachedExists {
		return cachedToxicity, nil
	}

	var response float64

	requestUrl, err := url.Parse(c.config.ToxicityAddr)
	if err != nil {
		return false, fmt.Errorf("cannot parse toxicity addr: %w", err)
	}

	query := url.Values{}
	query.Set("text", text)
	requestUrl.RawQuery = query.Encode()

	resp, err := req.R().
		SetContext(ctx).
		SetSuccessResult(&response).
		Get(requestUrl.String())
	if err != nil {
		return false, fmt.Errorf("cannot request toxicity: %w", err)
	}
	if !resp.IsSuccessState() {
		return false, fmt.Errorf("cannot request toxicity: %s", resp.String())
	}

	isToxic := response == 1
	if err := c.redisClient.Set(ctx, redisKey, isToxic, 24*7*time.Hour).Err(); err != nil {
		return false, fmt.Errorf("cannot set cached toxicity: %w", err)
	}

	return isToxic, nil
}

func (c *Service) CheckTextsToxicity(ctx context.Context, texts []string) ([]bool, error) {
	wg, wgCtx := errgroup.WithContext(ctx)

	toxicities := make([]bool, len(texts))
	var toxicitiesMu sync.Mutex

	for idx, text := range texts {
		wg.Go(
			func() error {
				result, err := c.CheckTextToxicity(wgCtx, text)
				if err != nil {
					return err
				}

				toxicitiesMu.Lock()
				toxicities[idx] = result
				toxicitiesMu.Unlock()

				return nil
			},
		)
	}

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	return toxicities, nil
}
