package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type ChannelsCommandsResponses struct {
	ID                string         `gorm:"primaryKey;column:id;type:TEXT;default:uuid_generate_v4()" json:"id"`
	Text              null.String    `gorm:"column:text;type:TEXT;"                         json:"text"      swaggertype:"string"`
	CommandID         string         `gorm:"column:commandId;type:TEXT;"                    json:"commandId"`
	Order             int            `gorm:"column:order;type:INT"                          json:"order"`
	TwitchCategoryIDs pq.StringArray `gorm:"column:twitch_category_id;type:UUID;" json:"twitchCategoryID"`
	OnlineOnly        bool           `gorm:"column:online_only;type:BOOL;"                 json:"onlineOnly"`
	OfflineOnly       bool           `gorm:"column:offline_only;type:BOOL;"                json:"offlineOnly"`
}

func (c *ChannelsCommandsResponses) TableName() string {
	return "channels_commands_responses"
}
