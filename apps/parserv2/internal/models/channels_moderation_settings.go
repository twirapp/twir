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


Table: channels_moderation_settings
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] type                                           USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
[ 2] channelId                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] enabled                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
[ 4] subscribers                                    BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
[ 5] vips                                           BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
[ 6] banTime                                        INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [600]
[ 7] banMessage                                     TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 8] warningMessage                                 TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 9] checkClips                                     BOOL                 null: true   primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
[10] triggerLength                                  INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [300]
[11] maxPercentage                                  INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [50]
[12] blackListSentences                             JSONB                null: true   primary: false  isArray: false  auto: false  col: JSONB           len: -1      default: [[]]


JSON Sample
-------------------------------------
{    "id": "OHqVUDedckWrTVustYNhuIyuL",    "type": "hvASLXGHgQhSmrYQFuRyFdKrj",    "channel_id": "qRsgdoajydfaMfAQaNdUIQsfG",    "enabled": true,    "subscribers": true,    "vips": false,    "ban_time": 50,    "ban_message": "iGUKwWgoEWWVceIlTPRnZtBbQ",    "warning_message": "pOBofcsgqRxMtrnpMllrBJHqp",    "check_clips": true,    "trigger_length": 43,    "max_percentage": 84,    "black_list_sentences": "ZtIWxwnJfPMpCtSqhmDmuIucL"}



*/

// ChannelsModerationSettings struct is a row record of the channels_moderation_settings table in the tsuwari database
type ChannelsModerationSettings struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;"  json:"id"`
	//[ 1] type                                           USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
	Type string `gorm:"column:type;type:VARCHAR;"                        json:"type"`
	//[ 2] channelId                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ChannelID string `gorm:"column:channelId;type:TEXT;"                      json:"channel_id"`
	//[ 3] enabled                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	Enabled bool `gorm:"column:enabled;type:BOOL;default:false;"          json:"enabled"`
	//[ 4] subscribers                                    BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	Subscribers bool `gorm:"column:subscribers;type:BOOL;default:false;"      json:"subscribers"`
	//[ 5] vips                                           BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	Vips bool `gorm:"column:vips;type:BOOL;default:false;"             json:"vips"`
	//[ 6] banTime                                        INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [600]
	BanTime int32 `gorm:"column:banTime;type:INT4;default:600;"            json:"ban_time"`
	//[ 7] banMessage                                     TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	BanMessage sql.NullString `gorm:"column:banMessage;type:TEXT;"                     json:"ban_message"`
	//[ 8] warningMessage                                 TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	WarningMessage sql.NullString `gorm:"column:warningMessage;type:TEXT;"                 json:"warning_message"`
	//[ 9] checkClips                                     BOOL                 null: true   primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	CheckClips sql.NullBool `gorm:"column:checkClips;type:BOOL;default:false;"       json:"check_clips"`
	//[10] triggerLength                                  INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [300]
	TriggerLength sql.NullInt64 `gorm:"column:triggerLength;type:INT4;default:300;"      json:"trigger_length"`
	//[11] maxPercentage                                  INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [50]
	MaxPercentage sql.NullInt64 `gorm:"column:maxPercentage;type:INT4;default:50;"       json:"max_percentage"`
	//[12] blackListSentences                             JSONB                null: true   primary: false  isArray: false  auto: false  col: JSONB           len: -1      default: [[]]
	BlackListSentences sql.NullString `gorm:"column:blackListSentences;type:JSONB;default:[];" json:"black_list_sentences"`
}

var channels_moderation_settingsTableInfo = &TableInfo{
	Name: "channels_moderation_settings",
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
			Name:               "type",
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
			GoFieldName:        "Type",
			GoFieldType:        "string",
			JSONFieldName:      "type",
			ProtobufFieldName:  "type",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		{
			Index:              2,
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
			ProtobufPos:        3,
		},

		{
			Index:              3,
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
			ProtobufPos:        4,
		},

		{
			Index:              4,
			Name:               "subscribers",
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
			GoFieldName:        "Subscribers",
			GoFieldType:        "bool",
			JSONFieldName:      "subscribers",
			ProtobufFieldName:  "subscribers",
			ProtobufType:       "bool",
			ProtobufPos:        5,
		},

		{
			Index:              5,
			Name:               "vips",
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
			GoFieldName:        "Vips",
			GoFieldType:        "bool",
			JSONFieldName:      "vips",
			ProtobufFieldName:  "vips",
			ProtobufType:       "bool",
			ProtobufPos:        6,
		},

		{
			Index:              6,
			Name:               "banTime",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "BanTime",
			GoFieldType:        "int32",
			JSONFieldName:      "ban_time",
			ProtobufFieldName:  "ban_time",
			ProtobufType:       "int32",
			ProtobufPos:        7,
		},

		{
			Index:              7,
			Name:               "banMessage",
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
			GoFieldName:        "BanMessage",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "ban_message",
			ProtobufFieldName:  "ban_message",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		{
			Index:              8,
			Name:               "warningMessage",
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
			GoFieldName:        "WarningMessage",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "warning_message",
			ProtobufFieldName:  "warning_message",
			ProtobufType:       "string",
			ProtobufPos:        9,
		},

		{
			Index:              9,
			Name:               "checkClips",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "BOOL",
			DatabaseTypePretty: "BOOL",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "BOOL",
			ColumnLength:       -1,
			GoFieldName:        "CheckClips",
			GoFieldType:        "sql.NullBool",
			JSONFieldName:      "check_clips",
			ProtobufFieldName:  "check_clips",
			ProtobufType:       "bool",
			ProtobufPos:        10,
		},

		{
			Index:              10,
			Name:               "triggerLength",
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
			GoFieldName:        "TriggerLength",
			GoFieldType:        "sql.NullInt64",
			JSONFieldName:      "trigger_length",
			ProtobufFieldName:  "trigger_length",
			ProtobufType:       "int32",
			ProtobufPos:        11,
		},

		{
			Index:              11,
			Name:               "maxPercentage",
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
			GoFieldName:        "MaxPercentage",
			GoFieldType:        "sql.NullInt64",
			JSONFieldName:      "max_percentage",
			ProtobufFieldName:  "max_percentage",
			ProtobufType:       "int32",
			ProtobufPos:        12,
		},

		{
			Index:              12,
			Name:               "blackListSentences",
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
			GoFieldName:        "BlackListSentences",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "black_list_sentences",
			ProtobufFieldName:  "black_list_sentences",
			ProtobufType:       "string",
			ProtobufPos:        13,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *ChannelsModerationSettings) TableName() string {
	return "channels_moderation_settings"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *ChannelsModerationSettings) BeforeSave(*gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *ChannelsModerationSettings) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *ChannelsModerationSettings) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *ChannelsModerationSettings) TableInfo() *TableInfo {
	return channels_moderation_settingsTableInfo
}
