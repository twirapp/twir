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


Table: notifications
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] imageSrc                                       TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] createdAt                                      TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [CURRENT_TIMESTAMP]
[ 3] userId                                         TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "DuIiohglsRVQIhgkKFFGbtNdw",    "image_src": "UAorPURptbNKQdQQUJEoGxmQt",    "created_at": "2215-01-06T23:22:41.53073624+03:00",    "user_id": "dTqGMFCvscREQKqVfOSdJlbtU"}



*/

// Notifications struct is a row record of the notifications table in the tsuwari database
type Notifications struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;"            json:"id"`
	//[ 1] imageSrc                                       TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ImageSrc sql.NullString `gorm:"column:imageSrc;type:TEXT;"                                 json:"image_src"`
	//[ 2] createdAt                                      TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [CURRENT_TIMESTAMP]
	CreatedAt time.Time `gorm:"column:createdAt;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"created_at"`
	//[ 3] userId                                         TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	UserID sql.NullString `gorm:"column:userId;type:TEXT;"                                   json:"user_id"`
}

var notificationsTableInfo = &TableInfo{
	Name: "notifications",
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
			Name:               "imageSrc",
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
			GoFieldName:        "ImageSrc",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "image_src",
			ProtobufFieldName:  "image_src",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		{
			Index:              2,
			Name:               "createdAt",
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
			GoFieldName:        "CreatedAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "created_at",
			ProtobufFieldName:  "created_at",
			ProtobufType:       "uint64",
			ProtobufPos:        3,
		},

		{
			Index:              3,
			Name:               "userId",
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
			GoFieldName:        "UserID",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "user_id",
			ProtobufFieldName:  "user_id",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (n *Notifications) TableName() string {
	return "notifications"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (n *Notifications) BeforeSave(*gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (n *Notifications) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (n *Notifications) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (n *Notifications) TableInfo() *TableInfo {
	return notificationsTableInfo
}
