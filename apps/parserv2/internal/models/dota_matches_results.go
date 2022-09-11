package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
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

Table: dota_matches_results
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] match_id                                       TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] players                                        JSONB                null: false  primary: false  isArray: false  auto: false  col: JSONB           len: -1      default: []
[ 3] radiant_win                                    BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: []
[ 4] game_mode                                      INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []

JSON Sample
-------------------------------------
{    "id": "nauQppXWcuwaltChZijFXnuGg",    "match_id": "WKFNatBXgYcFiFdYMhujxDAcZ",    "players": "PRwTsLhgZNRdFuawMUDILbgvm",    "radiant_win": false,    "game_mode": 69}
*/

// DotaMatchesResults struct is a row record of the dota_matches_results table in the tsuwari database
type DotaMatchesResults struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	//[ 1] match_id                                       TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	MatchID string `gorm:"column:match_id;type:TEXT;foreignKey:match_id"   json:"match_id"`
	//[ 2] players                                        JSONB                null: false  primary: false  isArray: false  auto: false  col: JSONB           len: -1      default: []
	Players string `gorm:"column:players;type:JSONB;"                      json:"players"`
	//[ 3] radiant_win                                    BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: []
	RadiantWin bool `gorm:"column:radiant_win;type:BOOL;"                   json:"radiant_win"`
	//[ 4] game_mode                                      INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
	GameMode int32 `gorm:"column:game_mode;type:INT4;"                     json:"game_mode"`
}

var dota_matches_resultsTableInfo = &TableInfo{
	Name: "dota_matches_results",
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
			Name:               "match_id",
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
			GoFieldName:        "MatchID",
			GoFieldType:        "string",
			JSONFieldName:      "match_id",
			ProtobufFieldName:  "match_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		{
			Index:              2,
			Name:               "players",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "JSONB",
			DatabaseTypePretty: "JSONB",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "JSONB",
			ColumnLength:       -1,
			GoFieldName:        "Players",
			GoFieldType:        "string",
			JSONFieldName:      "players",
			ProtobufFieldName:  "players",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		{
			Index:              3,
			Name:               "radiant_win",
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
			GoFieldName:        "RadiantWin",
			GoFieldType:        "bool",
			JSONFieldName:      "radiant_win",
			ProtobufFieldName:  "radiant_win",
			ProtobufType:       "bool",
			ProtobufPos:        4,
		},

		{
			Index:              4,
			Name:               "game_mode",
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
			GoFieldName:        "GameMode",
			GoFieldType:        "int32",
			JSONFieldName:      "game_mode",
			ProtobufFieldName:  "game_mode",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *DotaMatchesResults) TableName() string {
	return "dota_matches_results"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *DotaMatchesResults) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *DotaMatchesResults) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *DotaMatchesResults) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *DotaMatchesResults) TableInfo() *TableInfo {
	return dota_matches_resultsTableInfo
}
