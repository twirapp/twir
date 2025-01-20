package bots

const (
	SendMessageSubject   = "bots.send_message"
	DeleteMessageSubject = "bots.delete_message"
	BanSubject           = "bots.ban"
	ShoutOutSubject      = "bots.shoutout"
)

type SendMessageRequest struct {
	ChannelId         string
	ChannelName       *string
	Message           string
	IsAnnounce        bool
	SkipRateLimits    bool
	ReplyTo           string
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
	// BanTime set 0 to time permanent
	BanTime        int
	Reason         string
	IsModerator    bool
	AddModAfterBan bool
}

type SentShoutOutRequest struct {
	ChannelID string
	TargetID  string
}
