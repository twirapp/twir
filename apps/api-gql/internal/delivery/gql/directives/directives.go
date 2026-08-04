package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	model "github.com/twirapp/twir/libs/gomodels"
)

func New(sessions *auth.Auth, dashboardAccess *dashboardaccess.Service, rateLimiter *rate_limiter.LeakyBucketRateLimiter) *Directives {
	return &Directives{
		sessions:        sessions,
		dashboardAccess: dashboardAccess,
		rateLimiter:     rateLimiter,
	}
}

type Directives struct {
	sessions        sessionReader
	dashboardAccess *dashboardaccess.Service
	rateLimiter     *rate_limiter.LeakyBucketRateLimiter
}

type sessionReader interface {
	GetAuthenticatedUserModel(context.Context) (*model.Users, error)
	GetAuthenticatedUserByApiKey(context.Context) (*model.Users, error)
	GetSelectedDashboard(context.Context) (string, error)
}

func (c *Directives) NoRateLimit(
	ctx context.Context,
	obj any,
	next graphql.Resolver,
) (res any, err error) {
	return next(ctx)
}
