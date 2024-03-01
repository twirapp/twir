package bots

const (
	SendMessageSubject   = "bots.send_message"
	DeleteMessageSubject = "bots.delete_message"
)

type SendMessageRequest struct {
	ChannelId      string
	ChannelName    *string
	Message        string
	IsAnnounce     bool
	SkipRateLimits bool
	ReplyTo        string
}

type DeleteMessageRequest struct {
	ChannelId   string
	ChannelName *string
	MessageIds  []string
}
