package middlewares

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
)

func New(authService *auth.Auth, dashboardAccess *dashboardaccess.Service, humaAPI huma.API, rateLimiter *rate_limiter.LeakyBucketRateLimiter) *Middlewares {
	return &Middlewares{
		auth:            authService,
		dashboardAccess: dashboardAccess,
		huma:            humaAPI,
		rateLimiter:     rateLimiter,
	}
}

type Middlewares struct {
	auth            *auth.Auth
	dashboardAccess *dashboardaccess.Service
	huma            huma.API
	rateLimiter     *rate_limiter.LeakyBucketRateLimiter
}
