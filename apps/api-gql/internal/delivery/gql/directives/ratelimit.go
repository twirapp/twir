package directives

import (
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (c *Directives) RateLimit(
	ctx context.Context,
	obj any,
	next graphql.Resolver,
	key string,
	capacity int,
	window string,
) (res any, err error) {
	gCtx, err := gincontext.GetGinContext(ctx)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(window)
	if err != nil {
		return nil, fmt.Errorf("invalid duration format for rate limit window: %w", err)
	}

	resp, err := c.rateLimiter.Use(
		ctx,
		&rate_limiter.LeakyOptions{
			KeyPrefix:       fmt.Sprintf("api-gql:ratelimiter:%s:%s", gCtx.ClientIP(), key),
			MaximumCapacity: capacity,
			WindowSeconds:   int(duration.Seconds()),
		},
		1,
	)
	if err != nil {
		return nil, &gqlerror.Error{
			Err:     err,
			Message: "Rate limiter error",
		}
	}

	gCtx.Header("X-Rate-Limit-Bucket", key)
	gCtx.Header("X-Rate-Limit-Limit", fmt.Sprint(capacity))
	gCtx.Header("X-Rate-Limit-Remaining", fmt.Sprint(resp.RemainingTokens))
	gCtx.Header("X-Rate-Limit-Reset", fmt.Sprint(resp.ResetAt.Unix()))

	if !resp.Success {
		return nil, &gqlerror.Error{
			Message: "Rate limit exceeded",
			Err:     fmt.Errorf("rate limit exceeded"),
		}
	}

	return next(ctx)
}
