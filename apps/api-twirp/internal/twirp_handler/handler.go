package twirp_handler

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl"
	"github.com/satont/tsuwari/apps/api-twirp/internal/interceptors"
	"github.com/satont/tsuwari/apps/api-twirp/internal/wrappers"
	"github.com/satont/tsuwari/libs/grpc/generated/api"
	api_types "github.com/satont/tsuwari/libs/grpc/generated/api/api"
	"github.com/twitchtv/twirp"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Opts struct {
	Redis *redis.Client
	DB    *gorm.DB
}

func New(opts Opts) (string, http.Handler) {
	interceptorsService := interceptors.New(opts.Redis)

	twirpHandler := api.NewApiServer(
		impl.NewApi(impl.Opts{
			Redis: opts.Redis,
			DB:    opts.DB,
		}),
		twirp.WithServerPathPrefix("/v1"),
		twirp.WithServerInterceptors(interceptorsService.NewCacheInterceptor(interceptors.CacheOpts{
			CacheMethod:       "BotInfo",
			CacheDuration:     1 * time.Minute,
			WithChannelHeader: true,
			NewCastTo: func() any {
				return &api_types.BotInfo{}
			},
		})),
	)

	return twirpHandler.PathPrefix(), wrappers.Wrap(
		twirpHandler,
		wrappers.WithCors,
		wrappers.WithDashboardId,
	)
}
