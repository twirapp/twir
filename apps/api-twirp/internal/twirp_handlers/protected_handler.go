package twirp_handlers

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected"
	"github.com/satont/tsuwari/apps/api-twirp/internal/interceptors"
	"github.com/satont/tsuwari/apps/api-twirp/internal/wrappers"
	"github.com/satont/tsuwari/libs/grpc/generated/api"
	"github.com/satont/tsuwari/libs/grpc/generated/api/bots"
	"github.com/twitchtv/twirp"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Opts struct {
	Redis *redis.Client
	DB    *gorm.DB
}

func NewProtected(opts Opts) (string, http.Handler) {
	interceptorsService := interceptors.New(opts.Redis)

	twirpHandler := api.NewProtectedServer(
		impl_protected.New(impl_protected.Opts{
			Redis: opts.Redis,
			DB:    opts.DB,
		}),
		twirp.WithServerPathPrefix("/protected"),
		twirp.WithServerInterceptors(interceptorsService.NewCacheInterceptor(interceptors.CacheOpts{
			CacheMethod:       "BotInfo",
			CacheDuration:     1 * time.Minute,
			WithChannelHeader: true,
			NewCastTo: func() any {
				return &bots.BotInfo{}
			},
		})),
	)

	return twirpHandler.PathPrefix(), wrappers.Wrap(
		twirpHandler,
		wrappers.WithCors,
		wrappers.WithDashboardId,
	)
}
