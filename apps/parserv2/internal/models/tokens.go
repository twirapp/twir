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


Table: tokens
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] accessToken                                    TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] refreshToken                                   TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] expiresIn                                      INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
[ 4] obtainmentTimestamp                            TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "WryQbEmAVTAprqVfJwmRZwIov",    "access_token": "uxqIdDJurmJhyoggcxGfgcAAl",    "refresh_token": "qksbbthlfcoErkxTefWtHpeqx",    "expires_in": 86,    "obtainment_timestamp": "2077-08-30T06:07:52.590514424+03:00"}



*/

// Tokens struct is a row record of the tokens table in the tsuwari database
type Tokens struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	//[ 1] accessToken                                    TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	AccessToken string `gorm:"column:accessToken;type:TEXT;" json:"access_token"`
	//[ 2] refreshToken                                   TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	RefreshToken string `gorm:"column:refreshToken;type:TEXT;" json:"refresh_token"`
	//[ 3] expiresIn                                      INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
	ExpiresIn int32 `gorm:"column:expiresIn;type:INT4;" json:"expires_in"`
	//[ 4] obtainmentTimestamp                            TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
	ObtainmentTimestamp time.Time `gorm:"column:obtainmentTimestamp;type:TIMESTAMP;" json:"obtainment_timestamp"`
}

var tokensTableInfo = &TableInfo{
	Name: "tokens",
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
			Name:               "accessToken",
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
			GoFieldName:        "AccessToken",
			GoFieldType:        "string",
			JSONFieldName:      "access_token",
			ProtobufFieldName:  "access_token",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "refreshToken",
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
			GoFieldName:        "RefreshToken",
			GoFieldType:        "string",
			JSONFieldName:      "refresh_token",
			ProtobufFieldName:  "refresh_token",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "expiresIn",
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
			GoFieldName:        "ExpiresIn",
			GoFieldType:        "int32",
			JSONFieldName:      "expires_in",
			ProtobufFieldName:  "expires_in",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "obtainmentTimestamp",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TIMESTAMP",
			DatabaseTypePretty: "TIMESTAMP",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMP",
			ColumnLength:       -1,
			GoFieldName:        "ObtainmentTimestamp",
			GoFieldType:        "time.Time",
			JSONFieldName:      "obtainment_timestamp",
			ProtobufFieldName:  "obtainment_timestamp",
			ProtobufType:       "uint64",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (t *Tokens) TableName() string {
	return "tokens"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (t *Tokens) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (t *Tokens) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (t *Tokens) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (t *Tokens) TableInfo() *TableInfo {
	return tokensTableInfo
}
