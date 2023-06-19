package twirp_handlers

import (
	"github.com/satont/tsuwari/apps/api-twirp/internal/interceptors"
	"github.com/satont/tsuwari/apps/api-twirp/internal/wrappers"
	"github.com/satont/tsuwari/libs/grpc/generated/api"
	"github.com/satont/tsuwari/libs/grpc/generated/api/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/twitchtv/twirp"
	"go.uber.org/fx"
	"time"
)

type Grpc struct {
	fx.In

	Tokens tokens.TokensClient
}

func NewProtected(opts Opts) IHandler {
	twirpHandler := api.NewProtectedServer(
		opts.ImplProtected,
		twirp.WithServerPathPrefix("/v1"),
		twirp.WithServerInterceptors(opts.Interceptor.DbUserInterceptor),
		twirp.WithServerInterceptors(opts.Interceptor.ChannelAccessInterceptor),
		twirp.WithServerInterceptors(opts.Interceptor.NewCacheInterceptor(interceptors.CacheOpts{
			CacheMethod:       "BotInfo",
			CacheDuration:     5 * time.Second,
			WithChannelHeader: true,
			NewCastTo: func() any {
				return &bots.BotInfo{}
			},
		})),
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
