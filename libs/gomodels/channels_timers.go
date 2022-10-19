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

type ChannelsTimers struct {
	ID                       string                     `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;"      json:"id"`
	ChannelID                string                     `gorm:"column:channelId;type:TEXT;"                          json:"channel_id"`
	Name                     string                     `gorm:"column:name;type:VARCHAR;size:255;"                   json:"name"`
	Enabled                  bool                       `gorm:"column:enabled;type:BOOL;default:true;"               json:"enabled"`
	TimeInterval             int32                      `gorm:"column:timeInterval;type:INT4;default:0;"             json:"time_interval"`
	MessageInterval          int32                      `gorm:"column:messageInterval;type:INT4;default:0;"          json:"message_interval"`
	LastTriggerMessageNumber int32                      `gorm:"column:lastTriggerMessageNumber;type:INT4;default:0;" json:"last_trigger_message_number"`
	Responses                *[]ChannelsTimersResponses `gorm:"foreignKey:TimerID"                                   json:"responses"`
}

var channels_timersTableInfo = &TableInfo{
	Name: "channels_timers",
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
			Name:               "name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
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
			Name:               "responses",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "JSONB",
			DatabaseTypePretty: "JSONB",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "JSONB",
			ColumnLength:       -1,
			GoFieldName:        "Responses",
			GoFieldType:        "string",
			JSONFieldName:      "responses",
			ProtobufFieldName:  "responses",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		{
			Index:              5,
			Name:               "last",
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
			GoFieldName:        "Last",
			GoFieldType:        "int32",
			JSONFieldName:      "last",
			ProtobufFieldName:  "last",
			ProtobufType:       "int32",
			ProtobufPos:        6,
		},

		{
			Index:              6,
			Name:               "timeInterval",
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
			GoFieldName:        "TimeInterval",
			GoFieldType:        "int32",
			JSONFieldName:      "time_interval",
			ProtobufFieldName:  "time_interval",
			ProtobufType:       "int32",
			ProtobufPos:        7,
		},

		{
			Index:              7,
			Name:               "messageInterval",
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
			GoFieldName:        "MessageInterval",
			GoFieldType:        "int32",
			JSONFieldName:      "message_interval",
			ProtobufFieldName:  "message_interval",
			ProtobufType:       "int32",
			ProtobufPos:        8,
		},

		{
			Index:              8,
			Name:               "lastTriggerMessageNumber",
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
			GoFieldName:        "LastTriggerMessageNumber",
			GoFieldType:        "int32",
			JSONFieldName:      "last_trigger_message_number",
			ProtobufFieldName:  "last_trigger_message_number",
			ProtobufType:       "int32",
			ProtobufPos:        9,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *ChannelsTimers) TableName() string {
	return "channels_timers"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *ChannelsTimers) BeforeSave(*gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *ChannelsTimers) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *ChannelsTimers) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *ChannelsTimers) TableInfo() *TableInfo {
	return channels_timersTableInfo
}
