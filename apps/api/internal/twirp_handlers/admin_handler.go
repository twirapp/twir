package twirp_handlers

import (
	"github.com/bakins/twirpotel"
	"github.com/satont/twir/apps/api/internal/handlers"
	"github.com/satont/twir/apps/api/internal/hooks"
	"github.com/satont/twir/apps/api/internal/wrappers"
	"github.com/twirapp/twir/libs/api"
	"github.com/twitchtv/twirp"
)

func NewAdmin(opts Opts) handlers.IHandler {
	twirpHandler := api.NewAdminServer(
		opts.ImplAdmin,
		twirp.WithServerPathPrefix("/v1"),
		twirp.WithServerInterceptors(twirpotel.ServerInterceptor()),
		twirp.WithServerHooks(hooks.NewLoggingServerHooks(opts.Logger)),
		twirp.WithServerInterceptors(
			opts.Interceptor.DbUserInterceptor,
			opts.Interceptor.AdminInterceptor,
		),
	)

	h := handlers.New(
		handlers.Opts{
			Pattern: twirpHandler.PathPrefix(),
			Handler: wrappers.Wrap(
				twirpHandler,
				wrappers.WithCors,
				wrappers.WithApiKeyHeader,
			),
		},
	)

	return h
}
