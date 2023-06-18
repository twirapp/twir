package twirp_handlers

import (
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_unprotected"
	"github.com/satont/tsuwari/apps/api-twirp/internal/wrappers"
	"github.com/satont/tsuwari/libs/grpc/generated/api"
	"github.com/twitchtv/twirp"
	"net/http"
)

func NewUnProtected(opts Opts) (string, http.Handler) {
	twirpHandler := api.NewUnProtectedServer(
		impl_unprotected.New(impl_unprotected.Opts{
			Redis: opts.Redis,
			DB:    opts.DB,
		}),
		twirp.WithServerPathPrefix("/v1"),
	)

	return twirpHandler.PathPrefix(), wrappers.Wrap(
		twirpHandler,
		wrappers.WithCors,
		wrappers.WithDashboardId,
	)
}
