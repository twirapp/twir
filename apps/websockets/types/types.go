package types

type WebSocketMessage struct {
	EventName string `json:"eventName"`
	Data      any    `json:"data"`
	CreatedAt string `json:"createdAt"`
}
