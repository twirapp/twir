package twirp_handlers

import (
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_unprotected"
	"github.com/satont/tsuwari/apps/api-twirp/internal/interceptors"
	"go.uber.org/fx"
	"net/http"
)

type Opts struct {
	fx.In

	Interceptor     *interceptors.Service
	ImplProtected   *impl_protected.Protected
	ImplUnProtected *impl_unprotected.UnProtected
}

func AsHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(IHandler)),
		fx.ResultTags(`group:"handlers"`),
	)
}

type IHandler interface {
	Pattern() string
	Handler() http.Handler
}

type Handler struct {
	pattern string
	handler http.Handler
}

func (h *Handler) Pattern() string {
	return h.pattern
}

func (h *Handler) Handler() http.Handler {
	return h.handler
}
