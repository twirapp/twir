package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
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


Table: channels_keywords
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] channelId                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] text                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] response                                       TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 4] enabled                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [true]
[ 5] cooldown                                       INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]


JSON Sample
-------------------------------------
{    "id": "CwTXtQkayyrIauSFVsbhfZpdU",    "channel_id": "fGZpifWvnQFKhnNrwpIvkyYUt",    "text": "ZfpGhBubxFEEVREKOQSUGxyZZ",    "response": "SnPGVqndHXBFCJCALKEsjGpgP",    "enabled": false,    "cooldown": 23}



*/

// ChannelsKeywords struct is a row record of the channels_keywords table in the tsuwari database
type ChannelsKeywords struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	//[ 1] channelId                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ChannelID string `gorm:"column:channelId;type:TEXT;"                     json:"channel_id"`
	//[ 2] text                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Text string `gorm:"column:text;type:TEXT;"                          json:"text"`
	//[ 3] response                                       TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Response string `gorm:"column:response;type:TEXT;"                      json:"response"`
	//[ 4] enabled                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [true]
	Enabled bool `gorm:"column:enabled;type:BOOL;default:true;"          json:"enabled"`
	//[ 5] cooldown                                       INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]
	Cooldown sql.NullInt64 `gorm:"column:cooldown;type:INT4;default:0;"            json:"cooldown"`
}

var channels_keywordsTableInfo = &TableInfo{
	Name: "channels_keywords",
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
			ProtobufPos:        2,
		},

		{
			Index:              2,
			Name:               "text",
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
			GoFieldName:        "Text",
			GoFieldType:        "string",
			JSONFieldName:      "text",
			ProtobufFieldName:  "text",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		{
			Index:              3,
			Name:               "response",
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
			GoFieldName:        "Response",
			GoFieldType:        "string",
			JSONFieldName:      "response",
			ProtobufFieldName:  "response",
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
			ProtobufPos:        6,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *ChannelsKeywords) TableName() string {
	return "channels_keywords"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *ChannelsKeywords) BeforeSave(*gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *ChannelsKeywords) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *ChannelsKeywords) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *ChannelsKeywords) TableInfo() *TableInfo {
	return channels_keywordsTableInfo
}
