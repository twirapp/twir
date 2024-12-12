package entity

import (
	"time"
)

type DashboardWidgetEvent struct {
	ID        string
	ChannelID string
	UserID    string
	Type      DashboardWidgetEventType
	Data      DashboardWidgetEventData
	CreatedAt time.Time
}

type DashboardWidgetEventType string

func (e DashboardWidgetEventType) String() string {
	return string(e)
}

const (
	TypeDonation                   DashboardWidgetEventType = "DONATION"
	TypeFollow                     DashboardWidgetEventType = "FOLLOW"
	TypeRaided                     DashboardWidgetEventType = "RAIDED"
	TypeSubscribe                  DashboardWidgetEventType = "SUBSCRIBE"
	TypeReSubscribe                DashboardWidgetEventType = "RESUBSCRIBE"
	TypeSubGift                    DashboardWidgetEventType = "SUBGIFT"
	TypeFirstUserMessage           DashboardWidgetEventType = "FIRST_USER_MESSAGE"
	TypeChatClear                  DashboardWidgetEventType = "CHAT_CLEAR"
	TypeRedemptionCreated          DashboardWidgetEventType = "REDEMPTION_CREATED"
	TypeChannelBan                 DashboardWidgetEventType = "CHANNEL_BAN"
	TypeChannelUnbanRequestCreate  DashboardWidgetEventType = "CHANNEL_UNBAN_REQUEST_CREATE"
	TypeChannelUnbanRequestResolve DashboardWidgetEventType = "CHANNEL_UNBAN_REQUEST_RESOLVE"
)

type DashboardWidgetEventData struct {
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
