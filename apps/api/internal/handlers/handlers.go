package handlers

import (
	"net/http"

	"go.uber.org/fx"
)

type IHandler interface {
	Pattern() string
	Handler() http.Handler
}

func AsHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(IHandler)),
		fx.ResultTags(`group:"handlers"`),
	)
}

type Handler struct {
	pattern string
	handler http.Handler
}

type Opts struct {
	Pattern string
	Handler http.Handler
}

func New(opts Opts) *Handler {
	return &Handler{
		pattern: opts.Pattern,
		handler: opts.Handler,
	}
}

func (h *Handler) Pattern() string {
	return h.pattern
}

func (h *Handler) Handler() http.Handler {
	return h.handler
}
