package model

import (
	"time"
)

type Event struct {
	ID        string
	ChannelID string
	UserID    string
	Type      EventType
	Data      Data
	CreatedAt time.Time
}

type EventType string

func (e EventType) String() string {
	return string(e)
}

const (
	TypeDonation                   EventType = "DONATION"
	TypeFollow                     EventType = "FOLLOW"
	TypeRaided                     EventType = "RAIDED"
	TypeSubscribe                  EventType = "SUBSCRIBE"
	TypeReSubscribe                EventType = "RESUBSCRIBE"
	TypeSubGift                    EventType = "SUBGIFT"
	TypeFirstUserMessage           EventType = "FIRST_USER_MESSAGE"
	TypeChatClear                  EventType = "CHAT_CLEAR"
	TypeRedemptionCreated          EventType = "REDEMPTION_CREATED"
	TypeChannelBan                 EventType = "CHANNEL_BAN"
	TypeChannelUnbanRequestCreate  EventType = "CHANNEL_UNBAN_REQUEST_CREATE"
	TypeChannelUnbanRequestResolve EventType = "CHANNEL_UNBAN_REQUEST_RESOLVE"
)

type Data struct {
	//
	DonationAmount   string
	DonationCurrency string
	DonationMessage  string
	DonationUsername string

	//
	RaidedViewersCount    string
	RaidedFromUserName    string
	RaidedFromDisplayName string

	//
	FollowUserName        string
	FollowUserDisplayName string

	//
	RedemptionTitle           string
	RedemptionInput           string
	RedemptionUserName        string
	RedemptionUserDisplayName string
	RedemptionCost            string

	//
	SubLevel           string
	SubUserName        string
	SubUserDisplayName string

	//
	ReSubLevel           string
	ReSubUserName        string
	ReSubUserDisplayName string
	ReSubMonths          string
	ReSubStreak          string

	//
	SubGiftLevel                 string
	SubGiftUserName              string
	SubGiftUserDisplayName       string
	SubGiftTargetUserName        string
	SubGiftTargetUserDisplayName string

	//
	FirstUserMessageUserName        string
	FirstUserMessageUserDisplayName string
	FirstUserMessageMessage         string

	//
	ModeratorName        string
	ModeratorDisplayName string

	//
	BanReason        string
	BanEndsInMinutes string
	BannedUserName   string
	BannedUserLogin  string

	//
	UserLogin       string
	UserDisplayName string

	//
	Message string
}
