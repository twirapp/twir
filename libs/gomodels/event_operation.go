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

const (
	OperationBan               EventOperationType = "BAN"
	OperationUnban             EventOperationType = "UNBAN"
	OperationBanRandom         EventOperationType = "BAN_RANDOM"
	OperationVip               EventOperationType = "VIP"
	OperationUnvip             EventOperationType = "UNVIP"
	OperationUnvipRandom       EventOperationType = "UNVIP_RANDOM"
	OperationMod               EventOperationType = "MOD"
	OperationUnmod             EventOperationType = "UNMOD"
	OperationUnmodRandom       EventOperationType = "UNMOD_RANDOM"
	OperationSendMessage       EventOperationType = "SEND_MESSAGE"
	OperationChangeTitle       EventOperationType = "CHANGE_TITLE"
	OperationChangeCategory    EventOperationType = "CHANGE_CATEGORY"
	OperationFulfillRedemption EventOperationType = "FULFILL_REDEMPTION"
	OperationCancelRedemption  EventOperationType = "CANCEL_REDEMPTION"
	OperationEnableSubMode     EventOperationType = "ENABLE_SUBMODE"
	OperationDisableSubMode    EventOperationType = "DISABLE_SUBMODE"
	OperationEnableEmoteOnly   EventOperationType = "ENABLE_EMOTEONLY"
	OperationDisableEmoteOnly  EventOperationType = "DISABLE_EMOTEONLY"
)

type EventOperation struct {
	ID      string             `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	Type    EventOperationType `gorm:"column:type;type:TEXT;"                     json:"type"`
	Delay   null.Int           `gorm:"column:delay;type:int" json:"delay"`
	EventID string             `gorm:"column:eventId;type:string" json:"eventId"`

	Input  null.String `gorm:"column:input;type:string" json:"input"`
	Repeat int         `gorm:"column:repeat;type:int" json:"repeat"`
	Order  int         `gorm:"column:order;type:int" json:"order"`
}

func (c *EventOperation) TableName() string {
	return "channels_events_operations"
}
