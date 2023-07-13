package twirp_handlers

import (
	"github.com/satont/twir/apps/api-twirp/internal/handlers"
	"github.com/satont/twir/apps/api-twirp/internal/wrappers"
	"github.com/satont/twir/libs/grpc/generated/api"
	"github.com/twitchtv/twirp"
)

func NewUnProtected(opts Opts) handlers.IHandler {
	twirpHandler := api.NewUnProtectedServer(
		opts.ImplUnProtected,
		twirp.WithServerPathPrefix("/v1"),
		twirp.WithServerInterceptors(opts.Interceptor.Errors),
	)

	h := handlers.New(handlers.Opts{
		Pattern: twirpHandler.PathPrefix(),
		Handler: wrappers.Wrap(
			twirpHandler,
			wrappers.WithCors,
			wrappers.WithDashboardId,
			wrappers.WithApiKeyHeader,
		),
	})

	return h
}
