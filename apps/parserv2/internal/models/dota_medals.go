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


Table: dota_medals
[ 0] rank_tier                                      TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 1] name                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "rank_tier": "pVMrcUHwIgfGVycSdJLCVNQMG",    "name": "MsuiGLXoPpeGCDhQMAbsKRVur"}



*/

// DotaMedals struct is a row record of the dota_medals table in the tsuwari database
type DotaMedals struct {
	//[ 0] rank_tier                                      TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: []
	RankTier string `gorm:"primary_key;column:rank_tier;type:TEXT;" json:"rank_tier"`
	//[ 1] name                                           TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Name string `gorm:"column:name;type:TEXT;" json:"name"`
}

var dota_medalsTableInfo = &TableInfo{
	Name: "dota_medals",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "rank_tier",
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
			GoFieldName:        "RankTier",
			GoFieldType:        "string",
			JSONFieldName:      "rank_tier",
			ProtobufFieldName:  "rank_tier",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "name",
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
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *DotaMedals) TableName() string {
	return "dota_medals"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *DotaMedals) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *DotaMedals) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *DotaMedals) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *DotaMedals) TableInfo() *TableInfo {
	return dota_medalsTableInfo
}
