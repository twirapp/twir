package middlewares

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Auth           *auth.Auth
	ChannelService *channelservice.ChannelService
	Gorm           *gorm.DB
	Huma           huma.API
	RateLimiter    *rate_limiter.LeakyBucketRateLimiter
}

func New(opts Opts) *Middlewares {
	return &Middlewares{
		auth:          opts.Auth,
		channelGetter: opts.ChannelService,
		gorm:          opts.Gorm,
		huma:          opts.Huma,
		rateLimiter:   opts.RateLimiter,
	}
}

type Middlewares struct {
	auth          *auth.Auth
	channelGetter selectedDashboardChannelGetter
	gorm          *gorm.DB
	huma          huma.API
	rateLimiter   *rate_limiter.LeakyBucketRateLimiter
}
