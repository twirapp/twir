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


Table: notifications_messages
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] text                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] title                                          TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] langCode                                       USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
[ 4] notificationId                                 TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "fNMdxdBVYeAPNkfejVjPrynSS",    "text": "MlmxJlRFmIWLCaKmCtIpGkATS",    "title": "fOYTpocrRPSLKtvqmTntHmktZ",    "lang_code": "ugBMhMYBoOhcLFjcsHZfXaQIN",    "notification_id": "dVvClKAOYpObotvPFiubUKEgH"}



*/

// NotificationsMessages struct is a row record of the notifications_messages table in the tsuwari database
type NotificationsMessages struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	//[ 1] text                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Text string `gorm:"column:text;type:TEXT;"                          json:"text"`
	//[ 2] title                                          TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Title sql.NullString `gorm:"column:title;type:TEXT;"                         json:"title"`
	//[ 3] langCode                                       USER_DEFINED         null: false  primary: false  isArray: false  auto: false  col: USER_DEFINED    len: -1      default: []
	LangCode string `gorm:"column:langCode;type:VARCHAR;"                   json:"lang_code"`
	//[ 4] notificationId                                 TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	NotificationID string `gorm:"column:notificationId;type:TEXT;"                json:"notification_id"`
}

var notifications_messagesTableInfo = &TableInfo{
	Name: "notifications_messages",
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
			Name:               "text",
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
			GoFieldName:        "Text",
			GoFieldType:        "string",
			JSONFieldName:      "text",
			ProtobufFieldName:  "text",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		{
			Index:              2,
			Name:               "title",
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
			GoFieldName:        "Title",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "title",
			ProtobufFieldName:  "title",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		{
			Index:              3,
			Name:               "langCode",
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
			GoFieldName:        "LangCode",
			GoFieldType:        "string",
			JSONFieldName:      "lang_code",
			ProtobufFieldName:  "lang_code",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		{
			Index:              4,
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
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (n *NotificationsMessages) TableName() string {
	return "notifications_messages"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (n *NotificationsMessages) BeforeSave(*gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (n *NotificationsMessages) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (n *NotificationsMessages) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (n *NotificationsMessages) TableInfo() *TableInfo {
	return notifications_messagesTableInfo
}
