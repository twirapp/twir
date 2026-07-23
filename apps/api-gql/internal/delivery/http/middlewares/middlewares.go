package middlewares

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Auth            *auth.Auth
	DashboardAccess *dashboardaccess.Service
	Huma            huma.API
	RateLimiter     *rate_limiter.LeakyBucketRateLimiter
}

func New(opts Opts) *Middlewares {
	return &Middlewares{
		auth:            opts.Auth,
		dashboardAccess: opts.DashboardAccess,
		huma:            opts.Huma,
		rateLimiter:     opts.RateLimiter,
	}
}

type Middlewares struct {
	auth            *auth.Auth
	dashboardAccess *dashboardaccess.Service
	huma            huma.API
	rateLimiter     *rate_limiter.LeakyBucketRateLimiter
}
