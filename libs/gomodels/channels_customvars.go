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

type ChannelsCustomvars struct {
	ID          string      `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	Name        string      `gorm:"column:name;type:TEXT;"                          json:"name"`
	Description null.String `gorm:"column:description;type:TEXT;"                   json:"description" swaggertype:"string"`
	//[ 3] type                                           USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
	Type string `gorm:"column:type;type:VARCHAR;"                       json:"type"`
	//[ 4] evalValue                                      TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	EvalValue null.String `gorm:"column:evalValue;type:TEXT;"                     json:"evalValue"   swaggertype:"string"`
	//[ 5] response                                       TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Response null.String `gorm:"column:response;type:TEXT;"                      json:"response"    swaggertype:"string"`
	//[ 6] channelId                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ChannelID string `gorm:"column:channelId;type:TEXT;"                     json:"channelId"`
}

func (c *ChannelsCustomvars) TableName() string {
	return "channels_customvars"
}
