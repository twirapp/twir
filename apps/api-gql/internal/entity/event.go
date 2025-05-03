package entity

var EventNil = Event{}

type Event struct {
	ID          string
	ChannelID   string
	Type        EventType
	RewardID    *string
	CommandID   *string
	KeywordID   *string
	Description string
	Enabled     bool
	OnlineOnly  bool
	Operations  []EventOperation
}

type EventOperation struct {
	ID             string
	Type           EventOperationType
	Input          *string
	Delay          int
	Repeat         int
	UseAnnounce    bool
	TimeoutTime    int
	TimeoutMessage *string
	Target         *string
	Enabled        bool
	Filters        []EventOperationFilter
}

type EventOperationFilter struct {
	ID    string
	Type  string
	Left  string
	Right string
}

type EventType string

const (
	EventTypeFollow                     EventType = "FOLLOW"
	EventTypeSubscribe                  EventType = "SUBSCRIBE"
	EventTypeResubscribe                EventType = "RESUBSCRIBE"
	EventTypeSubGift                    EventType = "SUB_GIFT"
	EventTypeRedemptionCreated          EventType = "REDEMPTION_CREATED"
	EventTypeCommandUsed                EventType = "COMMAND_USED"
	EventTypeFirstUserMessage           EventType = "FIRST_USER_MESSAGE"
	EventTypeRaided                     EventType = "RAIDED"
	EventTypeTitleOrCategoryChanged     EventType = "TITLE_OR_CATEGORY_CHANGED"
	EventTypeStreamOnline               EventType = "STREAM_ONLINE"
	EventTypeStreamOffline              EventType = "STREAM_OFFLINE"
	EventTypeOnChatClear                EventType = "ON_CHAT_CLEAR"
	EventTypeDonate                     EventType = "DONATE"
	EventTypeKeywordMatched             EventType = "KEYWORD_MATCHED"
	EventTypeGreetingSended             EventType = "GREETING_SENDED"
	EventTypePollBegin                  EventType = "POLL_BEGIN"
	EventTypePollProgress               EventType = "POLL_PROGRESS"
	EventTypePollEnd                    EventType = "POLL_END"
	EventTypePredictionBegin            EventType = "PREDICTION_BEGIN"
	EventTypePredictionProgress         EventType = "PREDICTION_PROGRESS"
	EventTypePredictionEnd              EventType = "PREDICTION_END"
	EventTypePredictionLock             EventType = "PREDICTION_LOCK"
	EventTypeChannelBan                 EventType = "CHANNEL_BAN"
	EventTypeChannelUnbanRequestCreate  EventType = "CHANNEL_UNBAN_REQUEST_CREATE"
	EventTypeChannelUnbanRequestResolve EventType = "CHANNEL_UNBAN_REQUEST_RESOLVE"
	EventTypeChannelMessageDelete       EventType = "CHANNEL_MESSAGE_DELETE"
)

var AllEventType = []EventType{
	EventTypeFollow,
	EventTypeSubscribe,
	EventTypeResubscribe,
	EventTypeSubGift,
	EventTypeRedemptionCreated,
	EventTypeCommandUsed,
	EventTypeFirstUserMessage,
	EventTypeRaided,
	EventTypeTitleOrCategoryChanged,
	EventTypeStreamOnline,
	EventTypeStreamOffline,
	EventTypeOnChatClear,
	EventTypeDonate,
	EventTypeKeywordMatched,
	EventTypeGreetingSended,
	EventTypePollBegin,
	EventTypePollProgress,
	EventTypePollEnd,
	EventTypePredictionBegin,
	EventTypePredictionProgress,
	EventTypePredictionEnd,
	EventTypePredictionLock,
	EventTypeChannelBan,
	EventTypeChannelUnbanRequestCreate,
	EventTypeChannelUnbanRequestResolve,
	EventTypeChannelMessageDelete,
}

func (e EventType) IsValid() bool {
	switch e {
	case EventTypeFollow, EventTypeSubscribe, EventTypeResubscribe, EventTypeSubGift, EventTypeRedemptionCreated, EventTypeCommandUsed, EventTypeFirstUserMessage, EventTypeRaided, EventTypeTitleOrCategoryChanged, EventTypeStreamOnline, EventTypeStreamOffline, EventTypeOnChatClear, EventTypeDonate, EventTypeKeywordMatched, EventTypeGreetingSended, EventTypePollBegin, EventTypePollProgress, EventTypePollEnd, EventTypePredictionBegin, EventTypePredictionProgress, EventTypePredictionEnd, EventTypePredictionLock, EventTypeChannelBan, EventTypeChannelUnbanRequestCreate, EventTypeChannelUnbanRequestResolve, EventTypeChannelMessageDelete:
		return true
	}
	return false
}

func (e EventType) String() string {
	return string(e)
}

type EventOperationType string

const (
	EventOperationTypeSendMessage              EventOperationType = "SEND_MESSAGE"
	EventOperationTypeMessageDelete            EventOperationType = "MESSAGE_DELETE"
	EventOperationTypeTriggerAlert             EventOperationType = "TRIGGER_ALERT"
	EventOperationTypeTimeout                  EventOperationType = "TIMEOUT"
	EventOperationTypeTimeoutRandom            EventOperationType = "TIMEOUT_RANDOM"
	EventOperationTypeBan                      EventOperationType = "BAN"
	EventOperationTypeUnban                    EventOperationType = "UNBAN"
	EventOperationTypeBanRandom                EventOperationType = "BAN_RANDOM"
	EventOperationTypeVip                      EventOperationType = "VIP"
	EventOperationTypeUnvip                    EventOperationType = "UNVIP"
	EventOperationTypeUnvipRandom              EventOperationType = "UNVIP_RANDOM"
	EventOperationTypeUnvipRandomIfNoSlots     EventOperationType = "UNVIP_RANDOM_IF_NO_SLOTS"
	EventOperationTypeMod                      EventOperationType = "MOD"
	EventOperationTypeUnmod                    EventOperationType = "UNMOD"
	EventOperationTypeRaidChannel              EventOperationType = "RAID_CHANNEL"
	EventOperationTypeChangeVariable           EventOperationType = "CHANGE_VARIABLE"
	EventOperationTypeIncrementVariable        EventOperationType = "INCREMENT_VARIABLE"
	EventOperationTypeDecrementVariable        EventOperationType = "DECREMENT_VARIABLE"
	EventOperationTypeTtsSay                   EventOperationType = "TTS_SAY"
	EventOperationTypeTtsSkip                  EventOperationType = "TTS_SKIP"
	EventOperationTypeTtsEnable                EventOperationType = "TTS_ENABLE"
	EventOperationTypeTtsDisable               EventOperationType = "TTS_DISABLE"
	EventOperationTypeAllowCommandToUser       EventOperationType = "ALLOW_COMMAND_TO_USER"
	EventOperationTypeRemoveAllowCommandToUser EventOperationType = "REMOVE_ALLOW_COMMAND_TO_USER"
)

var AllEventOperationType = []EventOperationType{
	EventOperationTypeSendMessage,
	EventOperationTypeMessageDelete,
	EventOperationTypeTriggerAlert,
	EventOperationTypeTimeout,
	EventOperationTypeTimeoutRandom,
	EventOperationTypeBan,
	EventOperationTypeUnban,
	EventOperationTypeBanRandom,
	EventOperationTypeVip,
	EventOperationTypeUnvip,
	EventOperationTypeUnvipRandom,
	EventOperationTypeUnvipRandomIfNoSlots,
	EventOperationTypeMod,
	EventOperationTypeUnmod,
	EventOperationTypeRaidChannel,
	EventOperationTypeChangeVariable,
	EventOperationTypeIncrementVariable,
	EventOperationTypeDecrementVariable,
	EventOperationTypeTtsSay,
	EventOperationTypeTtsSkip,
	EventOperationTypeTtsEnable,
	EventOperationTypeTtsDisable,
	EventOperationTypeAllowCommandToUser,
	EventOperationTypeRemoveAllowCommandToUser,
}

func (e EventOperationType) IsValid() bool {
	switch e {
	case EventOperationTypeSendMessage, EventOperationTypeMessageDelete, EventOperationTypeTriggerAlert, EventOperationTypeTimeout, EventOperationTypeTimeoutRandom, EventOperationTypeBan, EventOperationTypeUnban, EventOperationTypeBanRandom, EventOperationTypeVip, EventOperationTypeUnvip, EventOperationTypeUnvipRandom, EventOperationTypeUnvipRandomIfNoSlots, EventOperationTypeMod, EventOperationTypeUnmod, EventOperationTypeRaidChannel, EventOperationTypeChangeVariable, EventOperationTypeIncrementVariable, EventOperationTypeDecrementVariable, EventOperationTypeTtsSay, EventOperationTypeTtsSkip, EventOperationTypeTtsEnable, EventOperationTypeTtsDisable, EventOperationTypeAllowCommandToUser, EventOperationTypeRemoveAllowCommandToUser:
		return true
	}
	return false
}

func (e EventOperationType) String() string {
	return string(e)
}
