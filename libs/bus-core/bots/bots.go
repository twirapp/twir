package bots

const (
	SendMessageSubject   = "bots.send_message"
	DeleteMessageSubject = "bots.delete_message"
	BanSubject           = "bots.ban"
	BanMultipleSubject   = "bots.ban_multiple"
	ShoutOutSubject      = "bots.shoutout"
	VipSubject           = "bots.vip"
	UnVipSubject         = "bots.unvip"
)

type SendMessageRequest struct {
	ChannelName       *string
	ChannelId         string
	Message           string
	ReplyTo           string
	IsAnnounce        bool
	SkipRateLimits    bool
	SkipToxicityCheck bool
	AnnounceColor
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

type VipRequest struct {
	ChannelID string
	TargetID  string
}

type UnVipRequest struct {
	ChannelID string
	TargetID  string
}

type AnnounceColor int

func (c AnnounceColor) String() string {
	return [...]string{"primary", "blue", "green", "orange", "purple"}[c]
}

const (
	AnnounceColorPrimary AnnounceColor = iota
	AnnounceColorBlue
	AnnounceColorGreen
	AnnounceColorOrange
	AnnounceColorPurple
)
