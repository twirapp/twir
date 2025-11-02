package middlewares

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Auth        *auth.Auth
	Gorm        *gorm.DB
	Huma        huma.API
	RateLimiter *rate_limiter.LeakyBucketRateLimiter
}

func New(opts Opts) *Middlewares {
	return &Middlewares{
		auth:        opts.Auth,
		gorm:        opts.Gorm,
		huma:        opts.Huma,
		rateLimiter: opts.RateLimiter,
	}
}

type Middlewares struct {
	auth        *auth.Auth
	gorm        *gorm.DB
	huma        huma.API
	rateLimiter *rate_limiter.LeakyBucketRateLimiter
}
