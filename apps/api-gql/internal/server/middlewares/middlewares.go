package middlewares

import (
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Sessions    *auth.Auth
	Logger      logger.Logger
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
	logger      logger.Logger
	rateLimiter *rate_limiter.LeakyBucketRateLimiter
}
