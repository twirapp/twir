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


Table: dota_matches
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] startedAt                                      TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[ 2] lobby_type                                     INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
[ 3] gameModeId                                     INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
[ 4] players                                        _INT4                null: true   primary: false  isArray: false  auto: false  col: _INT4           len: -1      default: []
[ 5] players_heroes                                 _INT4                null: true   primary: false  isArray: false  auto: false  col: _INT4           len: -1      default: []
[ 6] weekend_tourney_bracket_round                  TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 7] weekend_tourney_skill_level                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 8] match_id                                       TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 9] avarage_mmr                                    INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
[10] lobbyId                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[11] finished                                       BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]


JSON Sample
-------------------------------------
{    "id": "QWLOmsutBLqrIsciSIiKcidDS",    "started_at": "2157-06-03T20:22:57.821017349+03:00",    "lobby_type": 36,    "game_mode_id": 52,    "players": "sxgnIgIVndEqFuyFCcOVfRaMU",    "players_heroes": "fYXrsokJYRUtrMuLSScLKHGhl",    "weekend_tourney_bracket_round": "bKyKIxQUXtKFGFWKdnNhyoYno",    "weekend_tourney_skill_level": 4,    "match_id": "wQsqpuQOGhBUQXnmsZoYiOAHx",    "avarage_mmr": true}



*/

// DotaMatches struct is a row record of the dota_matches table in the tsuwari database
type DotaMatches struct {
	//[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: true   col: TEXT            len: -1      default: [gen_random_uuid()]
	ID string `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	//[ 1] startedAt                                      TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
	StartedAt time.Time `gorm:"column:startedAt;type:TIMESTAMP;" json:"started_at"`
	//[ 2] lobby_type                                     INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
	LobbyType int32 `gorm:"column:lobby_type;type:INT4;" json:"lobby_type"`
	//[ 3] gameModeId                                     INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
	GameModeID int32 `gorm:"column:gameModeId;type:INT4;" json:"game_mode_id"`
	//[ 6] weekend_tourney_bracket_round                  TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	WeekendTourneyBracketRound sql.NullString `gorm:"column:weekend_tourney_bracket_round;type:TEXT;" json:"weekend_tourney_bracket_round"`
	//[ 7] weekend_tourney_skill_level                    TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	WeekendTourneySkillLevel sql.NullString `gorm:"column:weekend_tourney_skill_level;type:TEXT;" json:"weekend_tourney_skill_level"`
	//[ 8] match_id                                       TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	MatchID string `gorm:"column:match_id;type:TEXT;" json:"match_id"`
	//[ 9] avarage_mmr                                    INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: []
	AvarageMmr int32 `gorm:"column:avarage_mmr;type:INT4;" json:"avarage_mmr"`
	//[10] lobbyId                                        TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	LobbyID string `gorm:"column:lobbyId;type:TEXT;" json:"lobby_id"`
	//[11] finished                                       BOOL                 null: false  primary: false  isArray: false  auto: false  col: BOOL            len: -1      default: [false]
	Finished bool `gorm:"column:finished;type:BOOL;default:false;" json:"finished"`
}

var dota_matchesTableInfo = &TableInfo{
	Name: "dota_matches",
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
			Name:               "startedAt",
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
			GoFieldName:        "StartedAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "started_at",
			ProtobufFieldName:  "started_at",
			ProtobufType:       "uint64",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "lobby_type",
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
			GoFieldName:        "LobbyType",
			GoFieldType:        "int32",
			JSONFieldName:      "lobby_type",
			ProtobufFieldName:  "lobby_type",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "gameModeId",
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
			GoFieldName:        "GameModeID",
			GoFieldType:        "int32",
			JSONFieldName:      "game_mode_id",
			ProtobufFieldName:  "game_mode_id",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "weekend_tourney_bracket_round",
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
			GoFieldName:        "WeekendTourneyBracketRound",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "weekend_tourney_bracket_round",
			ProtobufFieldName:  "weekend_tourney_bracket_round",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "weekend_tourney_skill_level",
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
			GoFieldName:        "WeekendTourneySkillLevel",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "weekend_tourney_skill_level",
			ProtobufFieldName:  "weekend_tourney_skill_level",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
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
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
			Name:               "avarage_mmr",
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
			GoFieldName:        "AvarageMmr",
			GoFieldType:        "int32",
			JSONFieldName:      "avarage_mmr",
			ProtobufFieldName:  "avarage_mmr",
			ProtobufType:       "int32",
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
			Name:               "lobbyId",
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
			GoFieldName:        "LobbyID",
			GoFieldType:        "string",
			JSONFieldName:      "lobby_id",
			ProtobufFieldName:  "lobby_id",
			ProtobufType:       "string",
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
			Name:               "finished",
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
			GoFieldName:        "Finished",
			GoFieldType:        "bool",
			JSONFieldName:      "finished",
			ProtobufFieldName:  "finished",
			ProtobufType:       "bool",
			ProtobufPos:        12,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *DotaMatches) TableName() string {
	return "dota_matches"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *DotaMatches) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *DotaMatches) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *DotaMatches) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *DotaMatches) TableInfo() *TableInfo {
	return dota_matchesTableInfo
}
