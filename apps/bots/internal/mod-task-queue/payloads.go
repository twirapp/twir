package mod_task_queue

type TaskModUserPayload struct {
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
}
