package twirp_handlers

import (
	"github.com/twirapp/twir/apps/api/internal/impl_protected"
	"github.com/twirapp/twir/apps/api/internal/impl_unprotected"
	"github.com/twirapp/twir/apps/api/internal/interceptors"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Logger          logger.Logger
	Interceptor     *interceptors.Service
	ImplProtected   *impl_protected.Protected
	ImplUnProtected *impl_unprotected.UnProtected
}
