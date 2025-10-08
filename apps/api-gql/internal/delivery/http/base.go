package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
)

type BaseOutputJson[T any] struct {
	Body BaseOutputBodyJson[T]
}

type BaseOutputBodyJson[T any] struct {
	Data T `json:"data"`
}

func CreateBaseOutputJson[T any](data T) *BaseOutputJson[T] {
	return &BaseOutputJson[T]{
		Body: BaseOutputBodyJson[T]{Data: data},
	}
}

type Route[Input any, Output any] interface {
	GetMeta() huma.Operation
	Handler(ctx context.Context, input Input) (Output, error)
	Register(api huma.API)
}

type registerRoute interface {
	GetMeta() huma.Operation
	Register(api huma.API)
}

func AsFxRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(registerRoute)),
		fx.ResultTags(`group:"huma-routes"`),
	)
}

type RegisterRoutesOpts struct {
	fx.In

	Api    huma.API
	Routes []registerRoute `group:"huma-routes"`
}

func RegisterRoutes(opts RegisterRoutesOpts) {
	for _, r := range opts.Routes {
		r.Register(opts.Api)
	}
}
