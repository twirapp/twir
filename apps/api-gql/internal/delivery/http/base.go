package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
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
