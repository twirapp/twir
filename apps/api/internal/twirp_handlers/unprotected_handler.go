package twirp_handlers

import (
	"github.com/satont/twir/apps/api/internal/handlers"
	"github.com/satont/twir/apps/api/internal/hooks"
	"github.com/satont/twir/apps/api/internal/wrappers"
	"github.com/satont/twir/libs/grpc/generated/api"
	"github.com/twitchtv/twirp"
)

func NewUnProtected(opts Opts) handlers.IHandler {
	twirpHandler := api.NewUnProtectedServer(
		opts.ImplUnProtected,
		twirp.WithServerPathPrefix("/v1"),
		twirp.WithServerHooks(hooks.NewLoggingServerHooks(opts.Logger)),
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
