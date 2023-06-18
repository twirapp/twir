package twirp_handler

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl"
	"github.com/satont/tsuwari/apps/api-twirp/internal/interceptors"
	"github.com/satont/tsuwari/apps/api-twirp/internal/wrappers"
	"github.com/satont/tsuwari/libs/grpc/generated/api"
	api_types "github.com/satont/tsuwari/libs/grpc/generated/api/api"
	"github.com/twitchtv/twirp"
	"net/http"
	"time"
)

func New(redisClient *redis.Client) (string, http.Handler) {
	interceptorsService := interceptors.New(redisClient)

	twirpHandler := api.NewApiServer(
		&impl.Api{},
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
