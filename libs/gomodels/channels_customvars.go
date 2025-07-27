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

type CustomVarType string

const (
	CustomVarScript        CustomVarType = "SCRIPT"
	CustomVarText          CustomVarType = "TEXT"
	CustomVarNumber        CustomVarType = "NUMBER"
	CustomVarChatChangable CustomVarType = "CHAT_CHANGABLE"
)

type ChannelsCustomvars struct {
	ID             string        `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	Name           string        `gorm:"column:name;type:TEXT;"                          json:"name"`
	Description    null.String   `gorm:"column:description;type:TEXT;"                   json:"description" swaggertype:"string"`
	Type           CustomVarType `gorm:"column:type;type:VARCHAR;"                       json:"type"`
	EvalValue      string        `gorm:"column:evalValue;type:TEXT;"                     json:"evalValue"   swaggertype:"string"`
	Response       string        `gorm:"column:response;type:TEXT;"                      json:"response"    swaggertype:"string"`
	ChannelID      string        `gorm:"column:channelId;type:TEXT;"                     json:"channelId"`
	ScriptLanguage string        `gorm:"column:script_language;type:VARCHAR;"       json:"scriptLanguage"`
}

func (c *ChannelsCustomvars) TableName() string {
	return "channels_customvars"
}
