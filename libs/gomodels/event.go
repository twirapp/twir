package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type EventType string

func (e EventType) String() string {
	return string(e)
}

const (
	EventTypeFollow                 EventType = "FOLLOW"
	EventTypeSubscribe                        = "SUBSCRIBE"
	EventTypeResubscribe                      = "RESUBSCRIBE"
	EventTypeSubGift                          = "SUB_GIFT"
	EventTypeRedemptionCreated                = "REDEMPTION_CREATED"
	EventTypeCommandUsed                      = "COMMAND_USED"
	EventTypeFirstUserMessage                 = "FIRST_USER_MESSAGE"
	EventTypeRaided                           = "RAIDED"
	EventTypeTitleOrCategoryChanged           = "TITLE_OR_CATEGORY_CHANGED"
	EventTypeStreamOnline                     = "STREAM_ONLINE"
	EventTypeStreamOffline                    = "STREAM_OFFLINE"
	EventTypeOnChatClear                      = "ON_CHAT_CLEAR"
	EventTypeDonate                           = "DONATE"
	EventTypeKeywordMatched                   = "KEYWORD_MATCHED"
	EventTypeGreetingSended                   = "GREETING_SENDED"
	EventTypePollBegin                        = "POLL_BEGIN"
	EventTypePollProgress                     = "POLL_PROGRESS"
	EventTypePollEnd                          = "POLL_END"
	EventTypePredictionBegin                  = "PREDICTION_BEGIN"
	EventTypePredictionProgress               = "PREDICTION_PROGRESS"
	EventTypePredictionEnd                    = "PREDICTION_END"
	EventTypePredictionLock                   = "PREDICTION_LOCK"
	EventStreamFirstUserJoin                  = "STREAM_FIRST_USER_JOIN"
)

type Event struct {
	ID          string      `gorm:"primaryKey;column:id;type:TEXT;default:uuid_generate_v4()" json:"id"`
	ChannelID   string      `gorm:"column:channelId;type:TEXT;" json:"channelId"`
	Type        EventType   `gorm:"column:type;type:TEXT;"                     json:"type"`
	RewardID    null.String `gorm:"column:rewardId;type:TEXT;"                     json:"rewardId"`
	CommandID   null.String `gorm:"column:commandId;type:TEXT;"                     json:"commandId"`
	KeywordID   null.String `gorm:"column:keywordId;type:TEXT;"                     json:"keywordId"`
	Description null.String `gorm:"column:description;type:TEXT" json:"description"`
	Enabled     bool        `gorm:"column:enabled;type:BOOL" json:"enabled"`
	OnlineOnly  bool        `gorm:"column:online_only;type:BOOL" json:"onlineOnly"`

	Operations []EventOperation `json:"operations"`
	Channel    *Channels        `gorm:"foreignKey:ChannelID" json:"channel"`
}

func (c *Event) TableName() string {
	return "channels_events"
}
