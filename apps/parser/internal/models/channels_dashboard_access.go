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


Table: channels_dashboard_access
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] channelId                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] userId                                         TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "HSNWZWuNMPxnPeatOUedNghpj",    "channel_id": "JYcnthZkdeXXGSQjbXgCeXWdi",    "user_id": "lVDEbdbCuyoYSyccNKWREvJfH"}



*/

// ChannelsDashboardAccess struct is a row record of the channels_dashboard_access table in the tsuwari database
type ChannelsDashboardAccess struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	//[ 1] channelId                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ChannelID string `gorm:"column:channelId;type:TEXT;"                     json:"channel_id"`
	//[ 2] userId                                         TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	UserID string `gorm:"column:userId;type:TEXT;"                        json:"user_id"`
}

var channels_dashboard_accessTableInfo = &TableInfo{
	Name: "channels_dashboard_access",
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
			Name:               "userId",
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
			GoFieldName:        "UserID",
			GoFieldType:        "string",
			JSONFieldName:      "user_id",
			ProtobufFieldName:  "user_id",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *ChannelsDashboardAccess) TableName() string {
	return "channels_dashboard_access"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *ChannelsDashboardAccess) BeforeSave(*gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *ChannelsDashboardAccess) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *ChannelsDashboardAccess) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *ChannelsDashboardAccess) TableInfo() *TableInfo {
	return channels_dashboard_accessTableInfo
}
