package twirp_handlers

import (
	"github.com/satont/tsuwari/apps/api-twirp/internal/wrappers"
	"github.com/satont/tsuwari/libs/grpc/generated/api"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/twitchtv/twirp"
	"go.uber.org/fx"
)

type Grpc struct {
	fx.In

	Tokens tokens.TokensClient
}

func NewProtected(opts Opts) *Handler {
	twirpHandler := api.NewProtectedServer(
		opts.ImplProtected,
		twirp.WithServerPathPrefix("/v1"),
		twirp.WithServerInterceptors(opts.Interceptor.SessionInterceptor),
		//twirp.WithServerInterceptors(interceptorsService.NewCacheInterceptor(interceptors.CacheOpts{
		//	CacheMethod:       "BotInfo",
		//	CacheDuration:     1 * time.Minute,
		//	WithChannelHeader: true,
		//	NewCastTo: func() any {
		//		return &bots.BotInfo{}
		//	},
		//})),
	)

	h := &Handler{
		pattern: twirpHandler.PathPrefix(),
		handler: wrappers.Wrap(
			twirpHandler,
			wrappers.WithCors,
			wrappers.WithDashboardId,
			wrappers.WithApiKeyHeader,
		),
	}

	return h
}
