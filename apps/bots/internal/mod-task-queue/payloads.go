package mod_task_queue

import "github.com/google/uuid"

type TaskModUserPayload struct {
	ChannelID    string    `json:"channel_id"`
	TwitchUserID uuid.UUID `json:"twitch_user_id"`
	UserID       string    `json:"user_id"`
}
