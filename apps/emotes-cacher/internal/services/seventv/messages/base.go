package messages

type BaseMessageWithoutData struct {
	Operation int `json:"op"`
	S         int `json:"s"`
}

type BaseMessage[T any] struct {
	Data      T   `json:"d"`
	Operation int `json:"op"`
	S         int `json:"s"`
}
