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

type ChannelsCommands struct {
	ID                 string                      `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;"  json:"id"`
	Name               string                      `gorm:"column:name;type:TEXT;"                           json:"name"`
	Cooldown           null.Int                    `gorm:"column:cooldown;type:INT4;default:0;"             json:"cooldown"           swaggertype:"integer"`
	CooldownType       string                      `gorm:"column:cooldownType;type:VARCHAR;default:GLOBAL;" json:"cooldownType"`
	Enabled            bool                        `gorm:"column:enabled;type:BOOL;"           json:"enabled"`
	Aliases            pq.StringArray              `gorm:"column:aliases;type:text[];default:[];"           json:"aliases"`
	Description        null.String                 `gorm:"column:description;type:TEXT;"                    json:"description"        swaggertype:"string"`
	Visible            bool                        `gorm:"column:visible;type:BOOL;"           json:"visible"`
	ChannelID          string                      `gorm:"column:channelId;type:TEXT;"                      json:"channelId"`
	Default            bool                        `gorm:"column:default;type:BOOL;default:false;"          json:"default"`
	DefaultName        null.String                 `gorm:"column:defaultName;type:TEXT;"                    json:"defaultName"        swaggertype:"string"`
	Module             string                      `gorm:"column:module;type:VARCHAR;default:CUSTOM;"       json:"module"`
	Responses          []ChannelsCommandsResponses `gorm:"foreignKey:CommandID"                             json:"responses"`
	IsReply            bool                        `gorm:"column:is_reply;type:BOOL;"           json:"isReply"`
	KeepResponsesOrder bool                        `gorm:"column:keepResponsesOrder;type:BOOL;" json:"keepResponsesOrder"`
	DeniedUsersIDS     pq.StringArray              `gorm:"column:deniedUsersIds;type:text[];default:[];"    json:"deniedUsersIds"`
	AllowedUsersIDS    pq.StringArray              `gorm:"column:allowedUsersIds;type:text[];default:[];"   json:"allowedUsersIds"`
	RolesIDS           pq.StringArray              `gorm:"column:rolesIds;type:text[];default:[];" json:"rolesIds"`

	GroupID null.String          `gorm:"column:groupId;type:UUID" json:"groupId"`
	Group   *ChannelCommandGroup `gorm:"foreignKey:GroupID" json:"group"`
}

func (c *ChannelsCommands) TableName() string {
	return "channels_commands"
}
