package twirp_handlers

import (
	"log/slog"

	"github.com/twirapp/twir/apps/api/internal/impl_protected"
	"github.com/twirapp/twir/apps/api/internal/impl_unprotected"
	"github.com/twirapp/twir/apps/api/internal/interceptors"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Logger          *slog.Logger
	Interceptor     *interceptors.Service
	ImplProtected   *impl_protected.Protected
	ImplUnProtected *impl_unprotected.UnProtected
}
