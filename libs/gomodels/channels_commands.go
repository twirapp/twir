package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

/*
DB Table Details
-------------------------------------
Table: channels_commands
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] name                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] cooldown                                       INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]
[ 3] cooldownType                                   USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: [GLOBAL]
[ 4] enabled                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [true]
[ 5] aliases                                        JSONB                null: true   primary: false  isArray: false  auto: false  col: JSONB           len: -1      default: [[]]
[ 6] description                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 7] visible                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [true]
[ 8] channelId                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 9] permission                                     USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
[10] default                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
[11] defaultName                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[12] module                                         USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: [CUSTOM]
JSON Sample
-------------------------------------
{    "id": "KgBDDcysaYZQqesDoXTtDnMYq",    "name": "AqCmFFUFcXbewxZltXJudmTRv",    "cooldown": 74,    "cooldown_type": "FTPPBGrueKbfSPZDEIKrpxqXT",    "enabled": true,    "aliases": "dDokZDGcWqNWOTpPvFpVKBpRd",    "description": "XJoOuXFCoEAvDouGhuhAWZNPJ",    "visible": true,    "channel_id": "adWSVvFXPchOZcXDWLelolCAE",    "permission": "JboyCeurieaMWJGuwBlFKaEYx",    "default": true,    "default_name": "VlbyNxwckpBlroYHbjnbduxVh",    "module": "aeCgoVJclnmKoOyoCIgmBXZCj"}
*/

// ChannelsCommands struct is a row record of the channels_commands table in the tsuwari database
type ChannelsCommands struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;"  json:"id"`
	//[ 1] name                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Name string `gorm:"column:name;type:TEXT;"                           json:"name"`
	//[ 2] cooldown                                       INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]
	Cooldown null.Int `gorm:"column:cooldown;type:INT4;default:0;"             json:"cooldown"`
	//[ 3] cooldownType                                   USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: [GLOBAL]
	CooldownType string `gorm:"column:cooldownType;type:VARCHAR;default:GLOBAL;" json:"cooldownType"`
	//[ 4] enabled                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [true]
	Enabled bool `gorm:"column:enabled;type:BOOL;default:true;"           json:"enabled"`
	//[ 5] aliases                                        JSONB                null: true   primary: false  isArray: false  auto: false  col: JSONB           len: -1      default: [[]]
	Aliases pq.StringArray `gorm:"column:aliases;type:text[];default:[];"            json:"aliases"`
	//[ 6] description                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Description null.String `gorm:"column:description;type:TEXT;"                    json:"description"`
	//[ 7] visible                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [true]
	Visible bool `gorm:"column:visible;type:BOOL;default:true;"           json:"visible"`
	//[ 8] channelId                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ChannelID string `gorm:"column:channelId;type:TEXT;"                      json:"channelId"`
	//[ 9] permission                                     USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
	Permission string `gorm:"column:permission;type:VARCHAR;"                  json:"permission"`
	//[10] default                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	Default bool `gorm:"column:default;type:BOOL;default:false;"          json:"default"`
	//[11] defaultName                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	DefaultName null.String `gorm:"column:defaultName;type:TEXT;"                    json:"defaultName"`
	//[12] module                                         USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: [CUSTOM]
	Module    string                      `gorm:"column:module;type:VARCHAR;default:CUSTOM;"       json:"module"`
	Responses []ChannelsCommandsResponses `gorm:"foreignKey:CommandID" json:"responses"`
	IsReply   bool                        `gorm:"column:is_reply";type:BOOL;default:true;" json:"isReply"`
}

var channels_commandsTableInfo = &TableInfo{
	Name: "channels_commands",
	Columns: []*ColumnInfo{
		{
			Index:              0,
			Name:               "id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       true,
			IsAutoIncrement:    true,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "ID",
			GoFieldType:        "string",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		{
			Index:              1,
			Name:               "name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		{
			Index:              2,
			Name:               "cooldown",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "Cooldown",
			GoFieldType:        "sql.NullInt64",
			JSONFieldName:      "cooldown",
			ProtobufFieldName:  "cooldown",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		{
			Index:              3,
			Name:               "cooldownType",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "USER_DEFINED",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "USER_DEFINED",
			ColumnLength:       -1,
			GoFieldName:        "CooldownType",
			GoFieldType:        "string",
			JSONFieldName:      "cooldown_type",
			ProtobufFieldName:  "cooldown_type",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		{
			Index:              4,
			Name:               "enabled",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "BOOL",
			DatabaseTypePretty: "BOOL",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "BOOL",
			ColumnLength:       -1,
			GoFieldName:        "Enabled",
			GoFieldType:        "bool",
			JSONFieldName:      "enabled",
			ProtobufFieldName:  "enabled",
			ProtobufType:       "bool",
			ProtobufPos:        5,
		},

		{
			Index:              5,
			Name:               "aliases",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "JSONB",
			DatabaseTypePretty: "JSONB",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "JSONB",
			ColumnLength:       -1,
			GoFieldName:        "Aliases",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "aliases",
			ProtobufFieldName:  "aliases",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		{
			Index:              6,
			Name:               "description",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "Description",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "description",
			ProtobufFieldName:  "description",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		{
			Index:              7,
			Name:               "visible",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "BOOL",
			DatabaseTypePretty: "BOOL",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "BOOL",
			ColumnLength:       -1,
			GoFieldName:        "Visible",
			GoFieldType:        "bool",
			JSONFieldName:      "visible",
			ProtobufFieldName:  "visible",
			ProtobufType:       "bool",
			ProtobufPos:        8,
		},

		{
			Index:              8,
			Name:               "channelId",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "ChannelID",
			GoFieldType:        "string",
			JSONFieldName:      "channel_id",
			ProtobufFieldName:  "channel_id",
			ProtobufType:       "string",
			ProtobufPos:        9,
		},

		{
			Index:              9,
			Name:               "permission",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "USER_DEFINED",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "USER_DEFINED",
			ColumnLength:       -1,
			GoFieldName:        "Permission",
			GoFieldType:        "string",
			JSONFieldName:      "permission",
			ProtobufFieldName:  "permission",
			ProtobufType:       "string",
			ProtobufPos:        10,
		},

		{
			Index:              10,
			Name:               "default",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "BOOL",
			DatabaseTypePretty: "BOOL",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "BOOL",
			ColumnLength:       -1,
			GoFieldName:        "Default",
			GoFieldType:        "bool",
			JSONFieldName:      "default",
			ProtobufFieldName:  "default",
			ProtobufType:       "bool",
			ProtobufPos:        11,
		},

		{
			Index:              11,
			Name:               "defaultName",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "DefaultName",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "default_name",
			ProtobufFieldName:  "default_name",
			ProtobufType:       "string",
			ProtobufPos:        12,
		},

		{
			Index:              12,
			Name:               "module",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "USER_DEFINED",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "USER_DEFINED",
			ColumnLength:       -1,
			GoFieldName:        "Module",
			GoFieldType:        "string",
			JSONFieldName:      "module",
			ProtobufFieldName:  "module",
			ProtobufType:       "string",
			ProtobufPos:        13,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *ChannelsCommands) TableName() string {
	return "channels_commands"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *ChannelsCommands) BeforeSave(*gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *ChannelsCommands) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *ChannelsCommands) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *ChannelsCommands) TableInfo() *TableInfo {
	return channels_commandsTableInfo
}
