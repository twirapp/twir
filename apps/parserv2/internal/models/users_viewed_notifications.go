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


Table: users_viewed_notifications
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] userId                                         TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] notificationId                                 TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] createdAt                                      TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [CURRENT_TIMESTAMP]


JSON Sample
-------------------------------------
{    "id": "eUUxGJHpAgkMdBnuxbwDManXd",    "user_id": "FZpCheavRWLOHyXrRkcxGUbxK",    "notification_id": "jMkakerscWFSenptjYbNXYeEK",    "created_at": "2283-03-01T06:33:15.332246152+03:00"}



*/

// UsersViewedNotifications struct is a row record of the users_viewed_notifications table in the tsuwari database
type UsersViewedNotifications struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;"            json:"id"`
	//[ 1] userId                                         TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	UserID string `gorm:"column:userId;type:TEXT;"                                   json:"user_id"`
	//[ 2] notificationId                                 TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	NotificationID string `gorm:"column:notificationId;type:TEXT;"                           json:"notification_id"`
	//[ 3] createdAt                                      TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [CURRENT_TIMESTAMP]
	CreatedAt time.Time `gorm:"column:createdAt;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"created_at"`
}

var users_viewed_notificationsTableInfo = &TableInfo{
	Name: "users_viewed_notifications",
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
			ProtobufPos:        2,
		},

		{
			Index:              2,
			Name:               "notificationId",
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
			GoFieldName:        "NotificationID",
			GoFieldType:        "string",
			JSONFieldName:      "notification_id",
			ProtobufFieldName:  "notification_id",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		{
			Index:              3,
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
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (u *UsersViewedNotifications) TableName() string {
	return "users_viewed_notifications"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (u *UsersViewedNotifications) BeforeSave(*gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (u *UsersViewedNotifications) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (u *UsersViewedNotifications) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (u *UsersViewedNotifications) TableInfo() *TableInfo {
	return users_viewed_notificationsTableInfo
}
