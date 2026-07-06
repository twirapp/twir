package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Sessions    *auth.Auth
	Gorm        *gorm.DB
	RateLimiter *rate_limiter.LeakyBucketRateLimiter
}

func New(opts Opts) *Directives {
	return &Directives{
		sessions:    opts.Sessions,
		gorm:        opts.Gorm,
		rateLimiter: opts.RateLimiter,
	}
}

type Directives struct {
	sessions    *auth.Auth
	gorm        *gorm.DB
	rateLimiter *rate_limiter.LeakyBucketRateLimiter
}

func (c *Directives) NoRateLimit(
	ctx context.Context,
	obj any,
	next graphql.Resolver,
) (res any, err error) {
	return next(ctx)
}
