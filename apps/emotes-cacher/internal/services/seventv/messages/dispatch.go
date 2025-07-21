package messages

type MessageDispatch[T any] struct {
	Type string `json:"type"`
	Body string `json:"body"`
}
