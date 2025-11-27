package middlewares

import (
	"log/slog"

	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Sessions    *auth.Auth
	Logger      *slog.Logger
	RateLimiter *rate_limiter.LeakyBucketRateLimiter
}

func New(opts Opts) *Middlewares {
	return &Middlewares{
		sessions:    opts.Sessions,
		logger:      opts.Logger,
		rateLimiter: opts.RateLimiter,
	}
}

type Middlewares struct {
	sessions    *auth.Auth
	logger      *slog.Logger
	rateLimiter *rate_limiter.LeakyBucketRateLimiter
}
