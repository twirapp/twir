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


Table: users
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] tokenId                                        TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] isTester                                       BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
[ 3] isBotAdmin                                     BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]


JSON Sample
-------------------------------------
{    "id": "kvnymClNqLFUClHRcxoZCuxEO",    "token_id": "JZYSvmuLNmMWBtmpOGtUexYZC",    "is_tester": false,    "is_bot_admin": true}



*/

// Users struct is a row record of the users table in the tsuwari database
type Users struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	ID string `gorm:"primary_key;column:id;type:TEXT;"           json:"id"`
	//[ 1] tokenId                                        TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	TokenID sql.NullString `gorm:"column:tokenId;type:TEXT;"                  json:"token_id"`
	//[ 2] isTester                                       BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	IsTester bool `gorm:"column:isTester;type:BOOL;default:false;"   json:"is_tester"`
	//[ 3] isBotAdmin                                     BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	IsBotAdmin bool `gorm:"column:isBotAdmin;type:BOOL;default:false;" json:"is_bot_admin"`
}

var usersTableInfo = &TableInfo{
	Name: "users",
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

		{
			Index:              1,
			Name:               "tokenId",
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
			GoFieldName:        "TokenID",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "token_id",
			ProtobufFieldName:  "token_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		{
			Index:              2,
			Name:               "isTester",
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
			GoFieldName:        "IsTester",
			GoFieldType:        "bool",
			JSONFieldName:      "is_tester",
			ProtobufFieldName:  "is_tester",
			ProtobufType:       "bool",
			ProtobufPos:        3,
		},

		{
			Index:              3,
			Name:               "isBotAdmin",
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
			GoFieldName:        "IsBotAdmin",
			GoFieldType:        "bool",
			JSONFieldName:      "is_bot_admin",
			ProtobufFieldName:  "is_bot_admin",
			ProtobufType:       "bool",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (u *Users) TableName() string {
	return "users"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (u *Users) BeforeSave(*gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (u *Users) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (u *Users) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (u *Users) TableInfo() *TableInfo {
	return usersTableInfo
}
