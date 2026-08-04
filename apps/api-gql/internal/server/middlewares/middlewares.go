package middlewares

import (
	"log/slog"

	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
)

func New(sessions *auth.Auth, logger *slog.Logger, rateLimiter *rate_limiter.LeakyBucketRateLimiter) *Middlewares {
	return &Middlewares{
		sessions:    sessions,
		logger:      logger,
		rateLimiter: rateLimiter,
	}
}

type Middlewares struct {
	sessions    *auth.Auth
	logger      *slog.Logger
	rateLimiter *rate_limiter.LeakyBucketRateLimiter
}

func (m *Middlewares) RateLimitInstance() *rate_limiter.LeakyBucketRateLimiter {
	return m.rateLimiter
}
