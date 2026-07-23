package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	model "github.com/twirapp/twir/libs/gomodels"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Sessions        *auth.Auth
	DashboardAccess *dashboardaccess.Service
	RateLimiter     *rate_limiter.LeakyBucketRateLimiter
}

func New(opts Opts) *Directives {
	return &Directives{
		sessions:        opts.Sessions,
		dashboardAccess: opts.DashboardAccess,
		rateLimiter:     opts.RateLimiter,
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
