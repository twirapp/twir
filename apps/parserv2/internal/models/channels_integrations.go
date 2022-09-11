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


Table: channels_integrations
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] enabled                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
[ 2] channelId                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] integrationId                                  TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 4] accessToken                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 5] refreshToken                                   TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 6] clientId                                       TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 7] clientSecret                                   TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 8] apiKey                                         TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 9] data                                           JSONB                null: true   primary: false  isArray: false  auto: false  col: JSONB           len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "CjHhapXelkulRGTDVlnddMRZh",    "enabled": true,    "channel_id": "JGUbvrhYyiXfswTHXBBweHBMW",    "integration_id": "TioVdsPWMTVUyglqWMFnFYfEh",    "access_token": "xDQSvHCGdWujlupWwrBDwEfKX",    "refresh_token": "uJgMFUNbKidhNpJTcBiDGsdJA",    "client_id": "TXHLsdEYHwYMiPJcWDrqTgWfh",    "client_secret": "nnGXLDSnjDCRVWFdOBsnMnPJT",    "api_key": "ZkvqhHvivQWEKXbdJOMUTCEjD",    "data": "oESTLAwUEvSftfJWCompIxaxo"}



*/

// ChannelsIntegrations struct is a row record of the channels_integrations table in the tsuwari database
type ChannelsIntegrations struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	//[ 1] enabled                                        BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	Enabled bool `gorm:"column:enabled;type:BOOL;default:false;"         json:"enabled"`
	//[ 2] channelId                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ChannelID string `gorm:"column:channelId;type:TEXT;"                     json:"channel_id"`
	//[ 3] integrationId                                  TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	IntegrationID string `gorm:"column:integrationId;type:TEXT;"                 json:"integration_id"`
	//[ 4] accessToken                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	AccessToken sql.NullString `gorm:"column:accessToken;type:TEXT;"                   json:"access_token"`
	//[ 5] refreshToken                                   TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	RefreshToken sql.NullString `gorm:"column:refreshToken;type:TEXT;"                  json:"refresh_token"`
	//[ 6] clientId                                       TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ClientID sql.NullString `gorm:"column:clientId;type:TEXT;"                      json:"client_id"`
	//[ 7] clientSecret                                   TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ClientSecret sql.NullString `gorm:"column:clientSecret;type:TEXT;"                  json:"client_secret"`
	//[ 8] apiKey                                         TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	APIKey sql.NullString `gorm:"column:apiKey;type:TEXT;"                        json:"api_key"`
	//[ 9] data                                           JSONB                null: true   primary: false  isArray: false  auto: false  col: JSONB           len: -1      default: []
	Data sql.NullString `gorm:"column:data;type:JSONB;"                         json:"data"`
}

var channels_integrationsTableInfo = &TableInfo{
	Name: "channels_integrations",
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
			Name:               "integrationId",
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
			GoFieldName:        "IntegrationID",
			GoFieldType:        "string",
			JSONFieldName:      "integration_id",
			ProtobufFieldName:  "integration_id",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		{
			Index:              4,
			Name:               "accessToken",
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
			GoFieldName:        "AccessToken",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "access_token",
			ProtobufFieldName:  "access_token",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		{
			Index:              5,
			Name:               "refreshToken",
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
			GoFieldName:        "RefreshToken",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "refresh_token",
			ProtobufFieldName:  "refresh_token",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		{
			Index:              6,
			Name:               "clientId",
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
			GoFieldName:        "ClientID",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "client_id",
			ProtobufFieldName:  "client_id",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		{
			Index:              7,
			Name:               "clientSecret",
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
			GoFieldName:        "ClientSecret",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "client_secret",
			ProtobufFieldName:  "client_secret",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		{
			Index:              8,
			Name:               "apiKey",
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
			GoFieldName:        "APIKey",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "api_key",
			ProtobufFieldName:  "api_key",
			ProtobufType:       "string",
			ProtobufPos:        9,
		},

		{
			Index:              9,
			Name:               "data",
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
			GoFieldName:        "Data",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "data",
			ProtobufFieldName:  "data",
			ProtobufType:       "string",
			ProtobufPos:        10,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *ChannelsIntegrations) TableName() string {
	return "channels_integrations"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *ChannelsIntegrations) BeforeSave(*gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *ChannelsIntegrations) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *ChannelsIntegrations) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *ChannelsIntegrations) TableInfo() *TableInfo {
	return channels_integrationsTableInfo
}
