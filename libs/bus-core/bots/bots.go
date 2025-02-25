package bots

const (
	SendMessageSubject   = "bots.send_message"
	DeleteMessageSubject = "bots.delete_message"
	BanSubject           = "bots.ban"
	ShoutOutSubject      = "bots.shoutout"
)

type SendMessageRequest struct {
	ChannelName       *string
	ChannelId         string
	Message           string
	ReplyTo           string
	IsAnnounce        bool
	SkipRateLimits    bool
	SkipToxicityCheck bool
}

type DeleteMessageRequest struct {
	ChannelId   string
	ChannelName *string
	MessageIds  []string
}

type BanRequest struct {
	ChannelID string
	UserID    string
	Reason    string
	// BanTime set 0 to time permanent
	BanTime        int
	IsModerator    bool
	AddModAfterBan bool
}

type SentShoutOutRequest struct {
	ChannelID string
	TargetID  string
}
