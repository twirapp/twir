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


Table: channels_dota_accounts
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] channelId                                      TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "lBasPwuLXQrDPoTnsMfrJFFjc",    "channel_id": "GejvUIvAoTSAHFYcjWVioSaIb"}



*/

// ChannelsDotaAccounts struct is a row record of the channels_dota_accounts table in the tsuwari database
type ChannelsDotaAccounts struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;" json:"id"`
	//[ 1] channelId                                      TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ChannelID string `gorm:"primary_key;column:channelId;type:TEXT;" json:"channel_id"`
}

var channels_dota_accountsTableInfo = &TableInfo{
	Name: "channels_dota_accounts",
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
			Name:               "channelId",
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
			GoFieldName:        "ChannelID",
			GoFieldType:        "string",
			JSONFieldName:      "channel_id",
			ProtobufFieldName:  "channel_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *ChannelsDotaAccounts) TableName() string {
	return "channels_dota_accounts"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *ChannelsDotaAccounts) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *ChannelsDotaAccounts) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *ChannelsDotaAccounts) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *ChannelsDotaAccounts) TableInfo() *TableInfo {
	return channels_dota_accountsTableInfo
}
