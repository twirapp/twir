package twirp_handlers

import (
	"github.com/satont/twir/apps/api/internal/impl_protected"
	"github.com/satont/twir/apps/api/internal/impl_unprotected"
	"github.com/satont/twir/apps/api/internal/interceptors"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Opts struct {
	fx.In

	Logger          *zap.Logger
	Interceptor     *interceptors.Service
	ImplProtected   *impl_protected.Protected
	ImplUnProtected *impl_unprotected.UnProtected
}
