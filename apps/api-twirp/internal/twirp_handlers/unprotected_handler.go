package twirp_handlers

import (
	"github.com/satont/tsuwari/apps/api-twirp/internal/wrappers"
	"github.com/satont/tsuwari/libs/grpc/generated/api"
	"github.com/twitchtv/twirp"
)

func NewUnProtected(opts Opts) *Handler {
	twirpHandler := api.NewUnProtectedServer(
		opts.ImplUnProtected,
		twirp.WithServerPathPrefix("/v1"),
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
