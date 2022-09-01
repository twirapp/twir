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


Table: integrations
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] service                                        USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
[ 2] accessToken                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] refreshToken                                   TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 4] clientId                                       TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 5] clientSecret                                   TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 6] apiKey                                         TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 7] redirectUrl                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "xGtJulvYTkmdOjGOVgNeVsRgu",    "service": "EsRVXWlHkCsAVYpxSHZaKJAfa",    "access_token": "HJHFSvuHURTVOoWmSxHvvAkBB",    "refresh_token": "AKwdJvEmXAhrcvdrqkHahZgry",    "client_id": "QkeswjtynRZCvAiKhYBdmOhJJ",    "client_secret": "hqbSPWFwMKOjFHcNSSVSvLvtU",    "api_key": "rWiHcSOnKZAxgyvyOXJQdfnIJ",    "redirect_url": "kBruEpbpuBEagIMPYGJOLMYrp"}



*/

// Integrations struct is a row record of the integrations table in the tsuwari database
type Integrations struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	//[ 1] service                                        USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
	Service string `gorm:"column:service;type:VARCHAR;" json:"service"`
	//[ 2] accessToken                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	AccessToken sql.NullString `gorm:"column:accessToken;type:TEXT;" json:"access_token"`
	//[ 3] refreshToken                                   TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	RefreshToken sql.NullString `gorm:"column:refreshToken;type:TEXT;" json:"refresh_token"`
	//[ 4] clientId                                       TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ClientID sql.NullString `gorm:"column:clientId;type:TEXT;" json:"client_id"`
	//[ 5] clientSecret                                   TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ClientSecret sql.NullString `gorm:"column:clientSecret;type:TEXT;" json:"client_secret"`
	//[ 6] apiKey                                         TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	APIKey sql.NullString `gorm:"column:apiKey;type:TEXT;" json:"api_key"`
	//[ 7] redirectUrl                                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	RedirectURL sql.NullString `gorm:"column:redirectUrl;type:TEXT;" json:"redirect_url"`
}

var integrationsTableInfo = &TableInfo{
	Name: "integrations",
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

		&ColumnInfo{
			Index:              1,
			Name:               "service",
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
			GoFieldName:        "Service",
			GoFieldType:        "string",
			JSONFieldName:      "service",
			ProtobufFieldName:  "service",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
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
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
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
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
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
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
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
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
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
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "redirectUrl",
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
			GoFieldName:        "RedirectURL",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "redirect_url",
			ProtobufFieldName:  "redirect_url",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},
	},
}

// TableName sets the insert table name for this struct type
func (i *Integrations) TableName() string {
	return "integrations"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (i *Integrations) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (i *Integrations) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (i *Integrations) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (i *Integrations) TableInfo() *TableInfo {
	return integrationsTableInfo
}
