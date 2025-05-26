package model

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type ChannelsEventsListItem struct {
	ID        uuid.UUID
	ChannelID string
	UserID    *string
	Type      ChannelEventListItemType
	Data      *ChannelsEventsListItemData
	CreatedAt time.Time
}

type ChannelEventListItemType string

func (e ChannelEventListItemType) String() string {
	return string(e)
}

const (
	ChannelEventListItemTypeDonation                   ChannelEventListItemType = "DONATION"
	ChannelEventListItemTypeFollow                     ChannelEventListItemType = "FOLLOW"
	ChannelEventListItemTypeRaided                     ChannelEventListItemType = "RAIDED"
	ChannelEventListItemTypeSubscribe                  ChannelEventListItemType = "SUBSCRIBE"
	ChannelEventListItemTypeReSubscribe                ChannelEventListItemType = "RESUBSCRIBE"
	ChannelEventListItemTypeSubGift                    ChannelEventListItemType = "SUBGIFT"
	ChannelEventListItemTypeFirstUserMessage           ChannelEventListItemType = "FIRST_USER_MESSAGE"
	ChannelEventListItemTypeChatClear                  ChannelEventListItemType = "CHAT_CLEAR"
	ChannelEventListItemTypeRedemptionCreated          ChannelEventListItemType = "REDEMPTION_CREATED"
	ChannelEventListItemTypeChannelBan                 ChannelEventListItemType = "CHANNEL_BAN"
	ChannelEventListItemTypeChannelUnbanRequestCreate  ChannelEventListItemType = "CHANNEL_UNBAN_REQUEST_CREATE"
	ChannelEventListItemTypeChannelUnbanRequestResolve ChannelEventListItemType = "CHANNEL_UNBAN_REQUEST_RESOLVE"
)

type ChannelsEventsListItemData struct {
	//
	DonationAmount   string `json:"donationAmount,omitempty"`
	DonationCurrency string `json:"donationCurrency,omitempty"`
	DonationMessage  string `json:"donationMessage,omitempty"`
	DonationUsername string `json:"donationUsername,omitempty"`

	//
	RaidedViewersCount    string `json:"raidedViewersCount,omitempty"`
	RaidedFromUserName    string `json:"raidedFromUserName,omitempty"`
	RaidedFromDisplayName string `json:"raidedFromDisplayName,omitempty"`

	//
	FollowUserName        string `json:"followUserName,omitempty"`
	FollowUserDisplayName string `json:"followUserDisplayName,omitempty"`

	//
	RedemptionTitle           string `json:"redemptionTitle,omitempty"`
	RedemptionInput           string `json:"redemptionInput,omitempty"`
	RedemptionUserName        string `json:"redemptionUserName,omitempty"`
	RedemptionUserDisplayName string `json:"redemptionUserDisplayName,omitempty"`
	RedemptionCost            string `json:"redemptionCost,omitempty"`

	//
	SubLevel           string `json:"subLevel,omitempty"`
	SubUserName        string `json:"subUserName,omitempty"`
	SubUserDisplayName string `json:"subUserDisplayName,omitempty"`

	//
	ReSubLevel           string `json:"reSubLevel,omitempty"`
	ReSubUserName        string `json:"reSubUserName,omitempty"`
	ReSubUserDisplayName string `json:"reSubUserDisplayName,omitempty"`
	ReSubMonths          string `json:"reSubMonths,omitempty"`
	ReSubStreak          string `json:"reSubStreak,omitempty"`

	//
	SubGiftLevel                 string `json:"subGiftLevel,omitempty"`
	SubGiftUserName              string `json:"subGiftUserName,omitempty"`
	SubGiftUserDisplayName       string `json:"subGiftUserDisplayName,omitempty"`
	SubGiftTargetUserName        string `json:"subGiftTargetUserName,omitempty"`
	SubGiftTargetUserDisplayName string `json:"subGiftTargetUserDisplayName,omitempty"`

	//
	FirstUserMessageUserName        string `json:"firstUserMessageUserName,omitempty"`
	FirstUserMessageUserDisplayName string `json:"firstUserMessageUserDisplayName,omitempty"`
	FirstUserMessageMessage         string `json:"firstUserMessageMessage,omitempty"`

	//
	ModeratorName        string `json:"moderatorName,omitempty"`
	ModeratorDisplayName string `json:"moderatorDisplayName,omitempty"`

	//
	BanReason        string `json:"banReason,omitempty"`
	BanEndsInMinutes string `json:"banEndsInMinutes,omitempty"`
	BannedUserName   string `json:"bannedUserName,omitempty"`
	BannedUserLogin  string `json:"bannedUserLogin,omitempty"`

	//
	UserLogin       string `json:"userLogin,omitempty"`
	UserDisplayName string `json:"userDisplayName,omitempty"`

	//
	Message string `json:"message,omitempty"`
}

func (a ChannelsEventsListItemData) Value() (driver.Value, error) {
	return json.Marshal(a)
}

var ErrChannelsEventsListItemDataTypeAssertionFailed = errors.New("type assertion to []byte failed")

func (a *ChannelsEventsListItemData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return ErrChannelsEventsListItemDataTypeAssertionFailed
	}
	return json.Unmarshal(b, &a)
}
