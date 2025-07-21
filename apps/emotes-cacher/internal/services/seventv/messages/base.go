package messages

import (
	"github.com/satont/twir/apps/emotes-cacher/internal/services/seventv/operations"
)

type BaseMessageWithoutData struct {
	Operation operations.IncomingOp `json:"op"`
	S         int                   `json:"s"`
}

type BaseMessage[T any] struct {
	Data      T   `json:"d"`
	Operation int `json:"op"`
	S         int `json:"s"`
}
