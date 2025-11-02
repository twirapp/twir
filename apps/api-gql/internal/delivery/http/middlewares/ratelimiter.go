package middlewares

import (
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
)

func (c *Middlewares) RateLimit(key string, capacity int, window time.Duration) func(
	ctx huma.Context,
	next func(huma.Context),
) {
	return func(ctx huma.Context, next func(huma.Context)) {
		resp, err := c.rateLimiter.Use(
			ctx.Context(),
			&rate_limiter.LeakyOptions{
				KeyPrefix:       fmt.Sprintf("api-gql:ratelimiter:%s:%s", ctx.RemoteAddr(), key),
				MaximumCapacity: capacity,
				WindowSeconds:   int(window.Seconds()),
			},
			1,
		)
		if err != nil {
			huma.WriteErr(c.huma, ctx, 500, "Rate limiter error", err)
			return
		}

		ctx.SetHeader("X-Rate-Limit-Bucket", key)
		ctx.SetHeader("X-Rate-Limit-Limit", fmt.Sprint(capacity))
		ctx.SetHeader("X-Rate-Limit-Remaining", fmt.Sprint(resp.RemainingTokens))
		ctx.SetHeader("X-Rate-Limit-Reset", fmt.Sprint(resp.ResetAt.Unix()))

		if !resp.Success {
			huma.WriteErr(c.huma, ctx, 429, "Rate limit exceeded", nil)
			return
		}

		next(ctx)
	}
}
