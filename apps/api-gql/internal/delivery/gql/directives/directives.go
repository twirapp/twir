package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	model "github.com/twirapp/twir/libs/gomodels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Sessions       *auth.Auth
	Gorm           *gorm.DB
	ChannelService *channelservice.ChannelService
	RateLimiter    *rate_limiter.LeakyBucketRateLimiter
}

func New(opts Opts) *Directives {
	return &Directives{
		sessions:               opts.Sessions,
		channels:               opts.ChannelService,
		selectedDashboardStore: &gormSelectedDashboardStore{gorm: opts.Gorm},
		gorm:                   opts.Gorm,
		rateLimiter:            opts.RateLimiter,
	}
}

type Directives struct {
	sessions               sessionReader
	channels               channelReader
	selectedDashboardStore selectedDashboardStore
	gorm                   *gorm.DB
	rateLimiter            *rate_limiter.LeakyBucketRateLimiter
}

type sessionReader interface {
	GetAuthenticatedUserModel(context.Context) (*model.Users, error)
	GetAuthenticatedUserByApiKey(context.Context) (*model.Users, error)
	GetSelectedDashboard(context.Context) (string, error)
}

type channelReader interface {
	GetChannelByID(context.Context, uuid.UUID) (channelsmodel.Channel, error)
}

type selectedDashboardStore interface {
	GetSelectedDashboardChannel(context.Context, string) (model.Channels, error)
	GetSelectedDashboardRoles(context.Context, string, string) ([]model.ChannelRole, error)
	GetSelectedDashboardUserStat(context.Context, string, string) (model.UsersStats, error)
}

func (c *Directives) NoRateLimit(
	ctx context.Context,
	obj any,
	next graphql.Resolver,
) (res any, err error) {
	return next(ctx)
}
