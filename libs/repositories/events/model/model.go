package model

var Nil = Event{}

type Event struct {
	ID          string           `json:"id"`
	ChannelID   string           `json:"channelId"`
	Type        EventType        `json:"type"`
	RewardID    *string          `json:"rewardId"`
	CommandID   *string          `json:"commandId"`
	KeywordID   *string          `json:"keywordId"`
	Description string           `json:"description"`
	Enabled     bool             `json:"enabled"`
	OnlineOnly  bool             `json:"onlineOnly"`
	Operations  []EventOperation `json:"operations"`
}

type EventOperation struct {
	ID             string                 `json:"id"`
	Type           EventOperationType     `json:"type"`
	Input          *string                `json:"input"`
	Delay          int                    `json:"delay"`
	Repeat         int                    `json:"repeat"`
	UseAnnounce    bool                   `json:"useAnnounce"`
	TimeoutTime    int                    `json:"timeoutTime"`
	TimeoutMessage *string                `json:"timeoutMessage"`
	Target         *string                `json:"target"`
	Enabled        bool                   `json:"enabled"`
	Filters        []EventOperationFilter `json:"filters"`
}

type EventOperationFilter struct {
	ID    string                   `json:"id"`
	Type  EventOperationFilterType `json:"type"`
	Left  string                   `json:"left"`
	Right string                   `json:"right"`
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
	EventOperationTypeUnmodRandom              EventOperationType = "UNMOD_RANDOM"
	EventOperationTypeRaidChannel              EventOperationType = "RAID_CHANNEL"
	EventOperationTypeChangeVariable           EventOperationType = "CHANGE_VARIABLE"
	EventOperationTypeIncrementVariable        EventOperationType = "INCREMENT_VARIABLE"
	EventOperationTypeDecrementVariable        EventOperationType = "DECREMENT_VARIABLE"
	EventOperationTypeTtsSay                   EventOperationType = "TTS_SAY"
	EventOperationTypeTtsSkip                  EventOperationType = "TTS_SKIP"
	EventOperationTypeTtsEnable                EventOperationType = "TTS_ENABLE"
	EventOperationTypeTtsDisable               EventOperationType = "TTS_DISABLE"
	EventOperationTypeTtsSwitchAutoread        EventOperationType = "TTS_SWITCH_AUTOREAD"
	EventOperationTypeTtsEnableAutoread        EventOperationType = "TTS_ENABLE_AUTOREAD"
	EventOperationTypeTtsDisableAutoread       EventOperationType = "TTS_DISABLE_AUTOREAD"
	EventOperationTypeAllowCommandToUser       EventOperationType = "ALLOW_COMMAND_TO_USER"
	EventOperationTypeRemoveAllowCommandToUser EventOperationType = "REMOVE_ALLOW_COMMAND_TO_USER"
	EventOperationTypeDenyCommandToUser        EventOperationType = "DENY_COMMAND_TO_USER"
	EventOperationTypeRemoveDenyCommandToUser  EventOperationType = "REMOVE_DENY_COMMAND_TO_USER"
	EventOperationTypeChangeTitle              EventOperationType = "CHANGE_TITLE"
	EventOperationTypeChangeCategory           EventOperationType = "CHANGE_CATEGORY"
	EventOperationTypeEnableSubmode            EventOperationType = "ENABLE_SUBMODE"
	EventOperationTypeDisableSubmode           EventOperationType = "DISABLE_SUBMODE"
	EventOperationTypeEnableEmoteOnly          EventOperationType = "ENABLE_EMOTE_ONLY"
	EventOperationTypeDisableEmoteOnly         EventOperationType = "DISABLE_EMOTE_ONLY"
	EventOperationTypeCreateGreeting           EventOperationType = "CREATE_GREETING"
	EventOperationTypeObsChangeScene           EventOperationType = "OBS_CHANGE_SCENE"
	EventOperationTypeObsToggleSource          EventOperationType = "OBS_TOGGLE_SOURCE"
	EventOperationTypeObsToggleAudio           EventOperationType = "OBS_TOGGLE_AUDIO"
	EventOperationTypeObsSetAudioVolume        EventOperationType = "OBS_SET_AUDIO_VOLUME"
	EventOperationTypeObsDecreaseAudioVolume   EventOperationType = "OBS_DECREASE_AUDIO_VOLUME"
	EventOperationTypeObsIncreaseAudioVolume   EventOperationType = "OBS_INCREASE_AUDIO_VOLUME"
	EventOperationTypeObsEnableAudio           EventOperationType = "OBS_ENABLE_AUDIO"
	EventOperationTypeObsDisableAudio          EventOperationType = "OBS_DISABLE_AUDIO"
	EventOperationTypeObsStartStream           EventOperationType = "OBS_START_STREAM"
	EventOperationTypeObsStopStream            EventOperationType = "OBS_STOP_STREAM"
	EventOperationTypeSeventvAddEmote          EventOperationType = "SEVENTV_ADD_EMOTE"
	EventOperationTypeSeventvRemoveEmote       EventOperationType = "SEVENTV_REMOVE_EMOTE"
	EventOperationTypeShoutoutChannel          EventOperationType = "SHOUTOUT_CHANNEL"
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
	EventOperationTypeUnmodRandom,
	EventOperationTypeRaidChannel,
	EventOperationTypeChangeVariable,
	EventOperationTypeIncrementVariable,
	EventOperationTypeDecrementVariable,
	EventOperationTypeTtsSay,
	EventOperationTypeTtsSkip,
	EventOperationTypeTtsEnable,
	EventOperationTypeTtsDisable,
	EventOperationTypeTtsSwitchAutoread,
	EventOperationTypeTtsEnableAutoread,
	EventOperationTypeTtsDisableAutoread,
	EventOperationTypeAllowCommandToUser,
	EventOperationTypeRemoveAllowCommandToUser,
	EventOperationTypeDenyCommandToUser,
	EventOperationTypeRemoveDenyCommandToUser,
	EventOperationTypeChangeTitle,
	EventOperationTypeChangeCategory,
	EventOperationTypeEnableSubmode,
	EventOperationTypeDisableSubmode,
	EventOperationTypeEnableEmoteOnly,
	EventOperationTypeDisableEmoteOnly,
	EventOperationTypeCreateGreeting,
	EventOperationTypeObsChangeScene,
	EventOperationTypeObsToggleSource,
	EventOperationTypeObsToggleAudio,
	EventOperationTypeObsSetAudioVolume,
	EventOperationTypeObsDecreaseAudioVolume,
	EventOperationTypeObsIncreaseAudioVolume,
	EventOperationTypeObsEnableAudio,
	EventOperationTypeObsDisableAudio,
	EventOperationTypeObsStartStream,
	EventOperationTypeObsStopStream,
	EventOperationTypeSeventvAddEmote,
	EventOperationTypeSeventvRemoveEmote,
	EventOperationTypeShoutoutChannel,
}

func (e EventOperationType) IsValid() bool {
	switch e {
	case EventOperationTypeSendMessage, EventOperationTypeMessageDelete, EventOperationTypeTriggerAlert, EventOperationTypeTimeout, EventOperationTypeTimeoutRandom, EventOperationTypeBan, EventOperationTypeUnban, EventOperationTypeBanRandom, EventOperationTypeVip, EventOperationTypeUnvip, EventOperationTypeUnvipRandom, EventOperationTypeUnvipRandomIfNoSlots, EventOperationTypeMod, EventOperationTypeUnmod, EventOperationTypeUnmodRandom, EventOperationTypeRaidChannel, EventOperationTypeChangeVariable, EventOperationTypeIncrementVariable, EventOperationTypeDecrementVariable, EventOperationTypeTtsSay, EventOperationTypeTtsSkip, EventOperationTypeTtsEnable, EventOperationTypeTtsDisable, EventOperationTypeTtsSwitchAutoread, EventOperationTypeTtsEnableAutoread, EventOperationTypeTtsDisableAutoread, EventOperationTypeAllowCommandToUser, EventOperationTypeRemoveAllowCommandToUser, EventOperationTypeDenyCommandToUser, EventOperationTypeRemoveDenyCommandToUser, EventOperationTypeChangeTitle, EventOperationTypeChangeCategory, EventOperationTypeEnableSubmode, EventOperationTypeDisableSubmode, EventOperationTypeEnableEmoteOnly, EventOperationTypeDisableEmoteOnly, EventOperationTypeCreateGreeting, EventOperationTypeObsChangeScene, EventOperationTypeObsToggleSource, EventOperationTypeObsToggleAudio, EventOperationTypeObsSetAudioVolume, EventOperationTypeObsDecreaseAudioVolume, EventOperationTypeObsIncreaseAudioVolume, EventOperationTypeObsEnableAudio, EventOperationTypeObsDisableAudio, EventOperationTypeObsStartStream, EventOperationTypeObsStopStream, EventOperationTypeSeventvAddEmote, EventOperationTypeSeventvRemoveEmote, EventOperationTypeShoutoutChannel:
		return true
	}
	return false
}

func (e EventOperationType) String() string {
	return string(e)
}

type EventOperationFilterType string

const (
	EventOperationFilterTypeEquals              EventOperationFilterType = "EQUALS"
	EventOperationFilterTypeNotEquals           EventOperationFilterType = "NOT_EQUALS"
	EventOperationFilterTypeContains            EventOperationFilterType = "CONTAINS"
	EventOperationFilterTypeNotContains         EventOperationFilterType = "NOT_CONTAINS"
	EventOperationFilterTypeStartsWith          EventOperationFilterType = "STARTS_WITH"
	EventOperationFilterTypeEndsWith            EventOperationFilterType = "ENDS_WITH"
	EventOperationFilterTypeGreaterThan         EventOperationFilterType = "GREATER_THAN"
	EventOperationFilterTypeLessThan            EventOperationFilterType = "LESS_THAN"
	EventOperationFilterTypeGreaterThanOrEquals EventOperationFilterType = "GREATER_THAN_OR_EQUALS"
	EventOperationFilterTypeLessThanOrEquals    EventOperationFilterType = "LESS_THAN_OR_EQUALS"
	EventOperationFilterTypeIsEmpty             EventOperationFilterType = "IS_EMPTY"
	EventOperationFilterTypeIsNotEmpty          EventOperationFilterType = "IS_NOT_EMPTY"
)

func (f EventOperationFilterType) String() string {
	return string(f)
}

var AllEventOperationFilterType = []EventOperationFilterType{
	EventOperationFilterTypeEquals,
	EventOperationFilterTypeNotEquals,
	EventOperationFilterTypeContains,
	EventOperationFilterTypeNotContains,
	EventOperationFilterTypeStartsWith,
	EventOperationFilterTypeEndsWith,
	EventOperationFilterTypeGreaterThan,
	EventOperationFilterTypeLessThan,
	EventOperationFilterTypeGreaterThanOrEquals,
	EventOperationFilterTypeLessThanOrEquals,
	EventOperationFilterTypeIsEmpty,
	EventOperationFilterTypeIsNotEmpty,
}
