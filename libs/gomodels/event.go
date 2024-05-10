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
	EventTypeSubscribe              EventType = "SUBSCRIBE"
	EventTypeResubscribe            EventType = "RESUBSCRIBE"
	EventTypeSubGift                EventType = "SUB_GIFT"
	EventTypeRedemptionCreated      EventType = "REDEMPTION_CREATED"
	EventTypeCommandUsed            EventType = "COMMAND_USED"
	EventTypeFirstUserMessage       EventType = "FIRST_USER_MESSAGE"
	EventTypeRaided                 EventType = "RAIDED"
	EventTypeTitleOrCategoryChanged EventType = "TITLE_OR_CATEGORY_CHANGED"
	EventTypeStreamOnline           EventType = "STREAM_ONLINE"
	EventTypeStreamOffline          EventType = "STREAM_OFFLINE"
	EventTypeOnChatClear            EventType = "ON_CHAT_CLEAR"
	EventTypeDonate                 EventType = "DONATE"
	EventTypeKeywordMatched         EventType = "KEYWORD_MATCHED"
	EventTypeGreetingSended         EventType = "GREETING_SENDED"
	EventTypePollBegin              EventType = "POLL_BEGIN"
	EventTypePollProgress           EventType = "POLL_PROGRESS"
	EventTypePollEnd                EventType = "POLL_END"
	EventTypePredictionBegin        EventType = "PREDICTION_BEGIN"
	EventTypePredictionProgress     EventType = "PREDICTION_PROGRESS"
	EventTypePredictionEnd          EventType = "PREDICTION_END"
	EventTypePredictionLock         EventType = "PREDICTION_LOCK"
	EventStreamFirstUserJoin        EventType = "STREAM_FIRST_USER_JOIN"
	EventChannelBan                 EventType = "CHANNEL_BAN"
	EventChannelUnbanRequestCreate  EventType = "CHANNEL_UNBAN_REQUEST_CREATE"
	EventChannelUnbanRequestResolve EventType = "CHANNEL_UNBAN_REQUEST_RESOLVE"
	EventChannelMessageDelete       EventType = "CHANNEL_MESSAGE_DELETE"
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
