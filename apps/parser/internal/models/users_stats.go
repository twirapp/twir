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


Table: users_stats
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] userId                                         TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] channelId                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] messages                                       INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]
[ 4] watched                                        INT8                 null: false  primary: false  isArray: false  auto: false  col: INT8            len: -1      default: [0]


JSON Sample
-------------------------------------
{    "id": "KlOOZwPyvnbnYxjNCbGxfwlpq",    "user_id": "kmclpXhfvctXbSsiUiAMeeKmM",    "channel_id": "OejKKuTfVDwolwCUknXcQexMh",    "messages": 89,    "watched": 23}



*/

// UsersStats struct is a row record of the users_stats table in the tsuwari database
type UsersStats struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	//[ 1] userId                                         TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	UserID string `gorm:"column:userId;type:TEXT;"                        json:"user_id"`
	//[ 2] channelId                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	ChannelID string `gorm:"column:channelId;type:TEXT;"                     json:"channel_id"`
	//[ 3] messages                                       INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]
	Messages int32 `gorm:"column:messages;type:INT4;default:0;"            json:"messages"`
	//[ 4] watched                                        INT8                 null: false  primary: false  isArray: false  auto: false  col: INT8            len: -1      default: [0]
	Watched int64 `gorm:"column:watched;type:INT8;default:0;"             json:"watched"`
}

var users_statsTableInfo = &TableInfo{
	Name: "users_stats",
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
			Name:               "messages",
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
			GoFieldName:        "Messages",
			GoFieldType:        "int32",
			JSONFieldName:      "messages",
			ProtobufFieldName:  "messages",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		{
			Index:              4,
			Name:               "watched",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT8",
			DatabaseTypePretty: "INT8",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT8",
			ColumnLength:       -1,
			GoFieldName:        "Watched",
			GoFieldType:        "int64",
			JSONFieldName:      "watched",
			ProtobufFieldName:  "watched",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (u *UsersStats) TableName() string {
	return "users_stats"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (u *UsersStats) BeforeSave(*gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (u *UsersStats) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (u *UsersStats) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (u *UsersStats) TableInfo() *TableInfo {
	return users_statsTableInfo
}
