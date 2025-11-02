package rate_limiter

import (
	"context"
	"time"

	ratelimit "github.com/aidenwallis/go-ratelimiting/redis"
	adapter "github.com/aidenwallis/go-ratelimiting/redis/adapters/go-redis"
	"github.com/redis/go-redis/v9"
)

func NewLeakyBucket(r *redis.Client) *LeakyBucketRateLimiter {
	lb := ratelimit.NewLeakyBucket(adapter.NewAdapter(r))

	return &LeakyBucketRateLimiter{
		lb: lb,
	}
}

type LeakyBucketRateLimiter struct {
	lb ratelimit.LeakyBucket
}

type LeakyOptions struct {
	KeyPrefix       string
	MaximumCapacity int
	WindowSeconds   int
}

type LeakyResponse struct {
	Success         bool
	RemainingTokens int
	ResetAt         time.Time
}

func (c *LeakyBucketRateLimiter) Use(
	ctx context.Context,
	bucket *LeakyOptions,
	takeAmount int,
) (*LeakyResponse, error) {
	resp, err := c.lb.Use(
		ctx,
		&ratelimit.LeakyBucketOptions{
			KeyPrefix:       bucket.KeyPrefix,
			MaximumCapacity: bucket.MaximumCapacity,
			WindowSeconds:   bucket.WindowSeconds,
		},
		takeAmount,
	)
	if err != nil {
		return nil, err
	}

	return &LeakyResponse{
		Success:         resp.Success,
		RemainingTokens: resp.RemainingTokens,
		ResetAt:         resp.ResetAt,
	}, nil
}
