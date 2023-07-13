package twirp_handlers

import (
	"github.com/satont/twir/apps/api-twirp/internal/impl_protected"
	"github.com/satont/twir/apps/api-twirp/internal/impl_unprotected"
	"github.com/satont/twir/apps/api-twirp/internal/interceptors"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Interceptor     *interceptors.Service
	ImplProtected   *impl_protected.Protected
	ImplUnProtected *impl_unprotected.UnProtected
}
