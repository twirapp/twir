package twirp_handlers

import (
	"github.com/satont/tsuwari/apps/api-twirp/internal/wrappers"
	"github.com/satont/tsuwari/libs/grpc/generated/api"
	"github.com/twitchtv/twirp"
)

func NewUnProtected(opts Opts) IHandler {
	twirpHandler := api.NewUnProtectedServer(
		opts.ImplUnProtected,
		twirp.WithServerPathPrefix("/v1"),
		twirp.WithServerInterceptors(opts.Interceptor.Errors),
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
