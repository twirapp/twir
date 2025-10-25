package bots

import (
	"math/rand/v2"
)

const (
	SendMessageSubject     = "bots.send_message"
	DeleteMessageSubject   = "bots.delete_message"
	BanSubject             = "bots.ban"
	BanMultipleSubject     = "bots.ban_multiple"
	ShoutOutSubject        = "bots.shoutout"
	VipSubject             = "bots.vip"
	UnVipSubject           = "bots.unvip"
	ModeratorAddSubject    = "bots.moderator_add"
	ModeratorRemoveSubject = "bots.moderator_remove"
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

type ModeratorAddRequest struct {
	ChannelID string
	TargetID  string
}

type ModeratorRemoveRequest struct {
	ChannelID string
	TargetID  string
}

type AnnounceColor int

func (c AnnounceColor) String() string {
	return allAnnounceColors[c]
}

var allAnnounceColors = map[AnnounceColor]string{
	AnnounceColorPrimary: "primary",
	AnnounceColorBlue:    "blue",
	AnnounceColorGreen:   "green",
	AnnounceColorOrange:  "orange",
	AnnounceColorPurple:  "purple",
}

func RandomAnnounceColor() AnnounceColor {
	return AnnounceColor(rand.IntN(len(allAnnounceColors) - 1))
}

const (
	AnnounceColorRandom                = -1
	AnnounceColorPrimary AnnounceColor = iota - 1
	AnnounceColorBlue
	AnnounceColorGreen
	AnnounceColorOrange
	AnnounceColorPurple
)
