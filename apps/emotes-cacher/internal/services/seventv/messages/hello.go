package messages

type HelloMessage struct {
	HeartBeatInterval uint32 `json:"heartbeat_interval"`
	SessionID         string `json:"session_id"`
	SubscriptionLimit uint32 `json:"subscription_limit"`
}
