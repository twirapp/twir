package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
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


Table: channels
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] isEnabled                                      BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [true]
[ 2] isTwitchBanned                                 BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
[ 3] isBanned                                       BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
[ 4] botId                                          TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "FNCAQLyhFILDnujkJLaGLbovZ",    "is_enabled": false,    "is_twitch_banned": true,    "is_banned": false,    "bot_id": "OxqeyxKjQnDASSFYlxYmpCGhn"}



*/

// Channels struct is a row record of the channels table in the tsuwari database
type Channels struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;" json:"id"`
	//[ 1] isEnabled                                      BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [true]
	IsEnabled bool `gorm:"column:isEnabled;type:BOOL;default:true;" json:"is_enabled"`
	//[ 2] isTwitchBanned                                 BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	IsTwitchBanned bool `gorm:"column:isTwitchBanned;type:BOOL;default:false;" json:"is_twitch_banned"`
	//[ 3] isBanned                                       BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	IsBanned bool `gorm:"column:isBanned;type:BOOL;default:false;" json:"is_banned"`
	//[ 4] botId                                          TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	BotID string `gorm:"column:botId;type:TEXT;" json:"bot_id"`
}

var channelsTableInfo = &TableInfo{
	Name: "channels",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
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

		&ColumnInfo{
			Index:              1,
			Name:               "isEnabled",
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
			GoFieldName:        "IsEnabled",
			GoFieldType:        "bool",
			JSONFieldName:      "is_enabled",
			ProtobufFieldName:  "is_enabled",
			ProtobufType:       "bool",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "isTwitchBanned",
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
			GoFieldName:        "IsTwitchBanned",
			GoFieldType:        "bool",
			JSONFieldName:      "is_twitch_banned",
			ProtobufFieldName:  "is_twitch_banned",
			ProtobufType:       "bool",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "isBanned",
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
			GoFieldName:        "IsBanned",
			GoFieldType:        "bool",
			JSONFieldName:      "is_banned",
			ProtobufFieldName:  "is_banned",
			ProtobufType:       "bool",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "botId",
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
			GoFieldName:        "BotID",
			GoFieldType:        "string",
			JSONFieldName:      "bot_id",
			ProtobufFieldName:  "bot_id",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *Channels) TableName() string {
	return "channels"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *Channels) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *Channels) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *Channels) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *Channels) TableInfo() *TableInfo {
	return channelsTableInfo
}
