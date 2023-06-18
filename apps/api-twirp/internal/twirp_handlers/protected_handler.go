package twirp_handlers

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected"
	"github.com/satont/tsuwari/apps/api-twirp/internal/interceptors"
	"github.com/satont/tsuwari/apps/api-twirp/internal/wrappers"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/api"
	"github.com/twitchtv/twirp"
	"gorm.io/gorm"
	"net/http"
)

type Opts struct {
	Redis          *redis.Client
	DB             *gorm.DB
	SessionManager *scs.SessionManager
	Config         *config.Config
}

func NewProtected(opts Opts) (string, http.Handler) {
	interceptorsService := interceptors.New(opts.Redis, opts.SessionManager, opts.DB)

	twirpHandler := api.NewProtectedServer(
		impl_protected.New(impl_protected.Opts{
			Redis:          opts.Redis,
			DB:             opts.DB,
			Config:         opts.Config,
			SessionManager: opts.SessionManager,
		}),
		twirp.WithServerPathPrefix("/v1"),
		twirp.WithServerInterceptors(interceptorsService.SessionInterceptor),
		//twirp.WithServerInterceptors(interceptorsService.NewCacheInterceptor(interceptors.CacheOpts{
		//	CacheMethod:       "BotInfo",
		//	CacheDuration:     1 * time.Minute,
		//	WithChannelHeader: true,
		//	NewCastTo: func() any {
		//		return &bots.BotInfo{}
		//	},
		//})),
	)

	return twirpHandler.PathPrefix(), wrappers.Wrap(
		twirpHandler,
		wrappers.WithCors,
		wrappers.WithDashboardId,
		wrappers.WithApiKeyHeader,
	)
}
