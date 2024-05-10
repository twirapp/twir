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

type EventOperationType string

func (r EventOperationType) String() string {
	return string(r)
}

const (
	OperationTimeout                  EventOperationType = "TIMEOUT"
	OperationTimeoutRandom            EventOperationType = "TIMEOUT_RANDOM"
	OperationBan                      EventOperationType = "BAN"
	OperationUnban                    EventOperationType = "UNBAN"
	OperationBanRandom                EventOperationType = "BAN_RANDOM"
	OperationVip                      EventOperationType = "VIP"
	OperationUnvip                    EventOperationType = "UNVIP"
	OperationUnvipRandom              EventOperationType = "UNVIP_RANDOM"
	OperationUnvipRandomIfNoSlots     EventOperationType = "UNVIP_RANDOM_IF_NO_SLOTS"
	OperationMod                      EventOperationType = "MOD"
	OperationUnmod                    EventOperationType = "UNMOD"
	OperationUnmodRandom              EventOperationType = "UNMOD_RANDOM"
	OperationSendMessage              EventOperationType = "SEND_MESSAGE"
	OperationChangeTitle              EventOperationType = "CHANGE_TITLE"
	OperationChangeCategory           EventOperationType = "CHANGE_CATEGORY"
	OperationFulfillRedemption        EventOperationType = "FULFILL_REDEMPTION"
	OperationCancelRedemption         EventOperationType = "CANCEL_REDEMPTION"
	OperationEnableSubMode            EventOperationType = "ENABLE_SUBMODE"
	OperationDisableSubMode           EventOperationType = "DISABLE_SUBMODE"
	OperationEnableEmoteOnly          EventOperationType = "ENABLE_EMOTEONLY"
	OperationDisableEmoteOnly         EventOperationType = "DISABLE_EMOTEONLY"
	OperationCreateGreeting           EventOperationType = "CREATE_GREETING"
	OperationObsSetScene              EventOperationType = "OBS_SET_SCENE"
	OperationObsToggleSource          EventOperationType = "OBS_TOGGLE_SOURCE"
	OperationObsToggleAudio           EventOperationType = "OBS_TOGGLE_AUDIO"
	OperationObsSetVolume             EventOperationType = "OBS_AUDIO_SET_VOLUME"
	OperationObsIncreaseVolume        EventOperationType = "OBS_AUDIO_INCREASE_VOLUME"
	OperationObsDecreaseVolume        EventOperationType = "OBS_AUDIO_DECREASE_VOLUME"
	OperationObsEnableAudio           EventOperationType = "OBS_ENABLE_AUDIO"
	OperationObsDisableAudio          EventOperationType = "OBS_DISABLE_AUDIO"
	OperationObsStopStream            EventOperationType = "OBS_STOP_STREAM"
	OperationObsStartStream           EventOperationType = "OBS_START_STREAM"
	OperationChangeVariable           EventOperationType = "CHANGE_VARIABLE"
	OperationIncrementVariable        EventOperationType = "INCREMENT_VARIABLE"
	OperationDecrementVariable        EventOperationType = "DECREMENT_VARIABLE"
	OperationTTSSay                   EventOperationType = "TTS_SAY"
	OperationTTSSkip                  EventOperationType = "TTS_SKIP"
	OperationTTSEnable                EventOperationType = "TTS_ENABLE"
	OperationTTSDisable               EventOperationType = "TTS_DISABLE"
	OperationAllowCommandToUser       EventOperationType = "ALLOW_COMMAND_TO_USER"
	OperationRemoveAllowCommandToUser EventOperationType = "REMOVE_ALLOW_COMMAND_TO_USER"
	OperationDenyCommandToUser        EventOperationType = "DENY_COMMAND_TO_USER"
	OperationRemoveDenyCommandToUser  EventOperationType = "REMOVE_DENY_COMMAND_TO_USER"
	OperationTTSSwitchAutoRead        EventOperationType = "TTS_SWITCH_AUTOREAD"
	OperationTTSEnableAutoRead        EventOperationType = "TTS_ENABLE_AUTOREAD"
	OperationTTSDisableAutoRead       EventOperationType = "TTS_DISABLE_AUTOREAD"
	OperationTriggerAlert             EventOperationType = "TRIGGER_ALERT"
	OperationSevenTvAddEmote          EventOperationType = "SEVENTV_ADD_EMOTE"
	OperationSevenTvRemoveEmote       EventOperationType = "SEVENTV_REMOVE_EMOTE"
	OperationRaidChannel              EventOperationType = "RAID_CHANNEL"
	OperationShoutoutChannel          EventOperationType = "SHOUTOUT_CHANNEL"
	OperationMessageDelete            EventOperationType = "MESSAGE_DELETE"
)

type EventOperation struct {
	ID      string             `gorm:"primaryKey;column:id;type:TEXT;default:uuid_generate_v4()" json:"id"`
	Type    EventOperationType `gorm:"column:type;type:TEXT;"                     json:"type"`
	Delay   int                `gorm:"column:delay;type:int" json:"delay"`
	EventID string             `gorm:"column:eventId;type:string" json:"eventId"`

	Input          null.String `gorm:"column:input;type:string" json:"input"`
	Repeat         int         `gorm:"column:repeat;type:int" json:"repeat"`
	Order          int         `gorm:"column:order;type:int" json:"order"`
	UseAnnounce    bool        `gorm:"column:useAnnounce;type:BOOL" json:"useAnnounce"`
	TimeoutTime    int         `gorm:"column:timeoutTime;type:int" json:"timeoutTime"`
	TimeoutMessage null.String `gorm:"column:timeoutMessage;type:text" json:"timeoutMessage"`
	Target         null.String `gorm:"column:target;type:string" json:"target"`
	Enabled        bool        `gorm:"column:enabled;type:bool" json:"enabled"`

	Filters []*EventOperationFilter `gorm:"foreignkey:OperationID" json:"filters"`
}

func (c *EventOperation) TableName() string {
	return "channels_events_operations"
}
