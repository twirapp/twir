package http

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
