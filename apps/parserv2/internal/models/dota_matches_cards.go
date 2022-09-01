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


Table: dota_matches_cards
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] match_id                                       TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] account_id                                     TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 3] rank_tier                                      INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
[ 4] leaderboard_rank                               INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "cmaERvQeFPbPXHdRlbXWOtgtx",    "match_id": "hPgFxvuEtvehhBrCXaJqrKXIx",    "account_id": "iaItcBicbZKZabBFVaWhjtBwR",    "rank_tier": 59,    "leaderboard_rank": 17}



*/

// DotaMatchesCards struct is a row record of the dota_matches_cards table in the tsuwari database
type DotaMatchesCards struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	//[ 1] match_id                                       TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	MatchID string `gorm:"column:match_id;type:TEXT;" json:"match_id"`
	//[ 2] account_id                                     TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	AccountID string `gorm:"column:account_id;type:TEXT;" json:"account_id"`
	//[ 3] rank_tier                                      INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
	RankTier sql.NullInt64 `gorm:"column:rank_tier;type:INT4;" json:"rank_tier"`
	//[ 4] leaderboard_rank                               INT4                 null: true   primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
	LeaderboardRank sql.NullInt64 `gorm:"column:leaderboard_rank;type:INT4;" json:"leaderboard_rank"`
}

var dota_matches_cardsTableInfo = &TableInfo{
	Name: "dota_matches_cards",
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

		&ColumnInfo{
			Index:              2,
			Name:               "account_id",
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
			GoFieldName:        "AccountID",
			GoFieldType:        "string",
			JSONFieldName:      "account_id",
			ProtobufFieldName:  "account_id",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "rank_tier",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "RankTier",
			GoFieldType:        "sql.NullInt64",
			JSONFieldName:      "rank_tier",
			ProtobufFieldName:  "rank_tier",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "leaderboard_rank",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "LeaderboardRank",
			GoFieldType:        "sql.NullInt64",
			JSONFieldName:      "leaderboard_rank",
			ProtobufFieldName:  "leaderboard_rank",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *DotaMatchesCards) TableName() string {
	return "dota_matches_cards"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *DotaMatchesCards) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *DotaMatchesCards) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *DotaMatchesCards) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *DotaMatchesCards) TableInfo() *TableInfo {
	return dota_matches_cardsTableInfo
}
