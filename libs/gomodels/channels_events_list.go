package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type ChannelEventListItemType string

func (e ChannelEventListItemType) String() string {
	return string(e)
}

const (
	ChannelEventListItemTypeDonation          ChannelEventListItemType = "DONATION"
	ChannelEventListItemTypeFollow                                     = "FOLLOW"
	ChannelEventListItemTypeRaided                                     = "RAIDED"
	ChannelEventListItemTypeSubscribe                                  = "SUBSCRIBE"
	ChannelEventListItemTypeReSubscribe                                = "RESUBSCRIBE"
	ChannelEventListItemTypeSubGift                                    = "SUBGIFT"
	ChannelEventListItemTypeFirstUserMessage                           = "FIRST_USER_MESSAGE"
	ChannelEventListItemTypeChatClear                                  = "CHAT_CLEAR"
	ChannelEventListItemTypeRedemptionCreated                          = "REDEMPTION_CREATED"
	ChannelEventListItemTypeChannelBan                                 = "CHANNEL_BAN"
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
	ModeratorName        string `json:"moderatorName"`
	ModeratorDisplayName string `json:"moderatorDisplayName"`

	//
	BanReason        string `json:"banReason"`
	BanEndsInMinutes string `json:"banEndsInMinutes"`
	BannedUserName   string `json:"bannedUserName"`
	BannedUserLogin  string `json:"bannedUserLogin"`
}

func (a ChannelsEventsListItemData) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *ChannelsEventsListItemData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

type ChannelsEventsListItem struct {
	ID        string                      `gorm:"primary_key;column:id;type:TEXT;"`
	ChannelID string                      `gorm:"column:channel_id;type:TEXT;"`
	UserID    string                      `gorm:"column:user_id;type:TEXT"`
	Type      ChannelEventListItemType    `gorm:"column:type;type:TEXT;"`
	Data      *ChannelsEventsListItemData `gorm:"column:data;type:JSONB;"`
	CreatedAt time.Time                   `gorm:"column:created_at;data:timestamp;"`

	Channel *Channels `gorm:"foreignKey:ChannelID"`
}

func (c *ChannelsEventsListItem) TableName() string {
	return "channels_events_list"
}
