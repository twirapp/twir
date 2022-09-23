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


Table: _prisma_migrations
[ 0] id                                             VARCHAR(36)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 36      default: []
[ 1] checksum                                       VARCHAR(64)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 64      default: []
[ 2] finished_at                                    TIMESTAMPTZ          null: true   primary: false  isArray: false  auto: false  col: TIMESTAMPTZ     len: -1      default: []
[ 3] migration_name                                 VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
[ 4] logs                                           TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 5] rolled_back_at                                 TIMESTAMPTZ          null: true   primary: false  isArray: false  auto: false  col: TIMESTAMPTZ     len: -1      default: []
[ 6] started_at                                     TIMESTAMPTZ          null: false  primary: false  isArray: false  auto: false  col: TIMESTAMPTZ     len: -1      default: [now()]
[ 7] applied_steps_count                            INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]


JSON Sample
-------------------------------------
{    "id": "CXprUxeGRlbfAdtCZHPRXPFjA",    "checksum": "wOpiFdKhrEkWxQdABBkCAeFGK",    "finished_at": "2129-05-26T03:01:45.798564643+03:00",    "migration_name": "xofHHOgBsLDhSBFHFqHMftoSt",    "logs": "eNNYJOKCXKZWiLGabNcOQwQGA",    "rolled_back_at": "2045-12-02T04:10:06.719580861+03:00",    "started_at": "2197-04-12T14:11:25.989309988+03:00",    "applied_steps_count": 59}



*/

// PrismaMigrations struct is a row record of the _prisma_migrations table in the tsuwari database
type PrismaMigrations struct {
	//[ 0] id                                             VARCHAR(36)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 36      default: []
	ID string `gorm:"primary_key;column:id;type:VARCHAR;size:36;"     json:"id"`
	//[ 1] checksum                                       VARCHAR(64)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 64      default: []
	Checksum string `gorm:"column:checksum;type:VARCHAR;size:64;"           json:"checksum"`
	//[ 2] finished_at                                    TIMESTAMPTZ          null: true   primary: false  isArray: false  auto: false  col: TIMESTAMPTZ     len: -1      default: []
	FinishedAt time.Time `gorm:"column:finished_at;type:TIMESTAMPTZ;"            json:"finished_at"`
	//[ 3] migration_name                                 VARCHAR(255)         null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 255     default: []
	MigrationName string `gorm:"column:migration_name;type:VARCHAR;size:255;"    json:"migration_name"`
	//[ 4] logs                                           TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Logs sql.NullString `gorm:"column:logs;type:TEXT;"                          json:"logs"`
	//[ 5] rolled_back_at                                 TIMESTAMPTZ          null: true   primary: false  isArray: false  auto: false  col: TIMESTAMPTZ     len: -1      default: []
	RolledBackAt time.Time `gorm:"column:rolled_back_at;type:TIMESTAMPTZ;"         json:"rolled_back_at"`
	//[ 6] started_at                                     TIMESTAMPTZ          null: false  primary: false  isArray: false  auto: false  col: TIMESTAMPTZ     len: -1      default: [now()]
	StartedAt time.Time `gorm:"column:started_at;type:TIMESTAMPTZ;"             json:"started_at"`
	//[ 7] applied_steps_count                            INT4                 null: false  primary: false  isArray: false  auto: false  col: INT4            len: -1      default: [0]
	AppliedStepsCount int32 `gorm:"column:applied_steps_count;type:INT4;default:0;" json:"applied_steps_count"`
}

var _prisma_migrationsTableInfo = &TableInfo{
	Name: "_prisma_migrations",
	Columns: []*ColumnInfo{
		{
			Index:              0,
			Name:               "id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(36)",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       36,
			GoFieldName:        "ID",
			GoFieldType:        "string",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		{
			Index:              1,
			Name:               "checksum",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(64)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       64,
			GoFieldName:        "Checksum",
			GoFieldType:        "string",
			JSONFieldName:      "checksum",
			ProtobufFieldName:  "checksum",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		{
			Index:              2,
			Name:               "finished_at",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TIMESTAMPTZ",
			DatabaseTypePretty: "TIMESTAMPTZ",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMPTZ",
			ColumnLength:       -1,
			GoFieldName:        "FinishedAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "finished_at",
			ProtobufFieldName:  "finished_at",
			ProtobufType:       "uint64",
			ProtobufPos:        3,
		},

		{
			Index:              3,
			Name:               "migration_name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       255,
			GoFieldName:        "MigrationName",
			GoFieldType:        "string",
			JSONFieldName:      "migration_name",
			ProtobufFieldName:  "migration_name",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		{
			Index:              4,
			Name:               "logs",
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
			GoFieldName:        "Logs",
			GoFieldType:        "sql.NullString",
			JSONFieldName:      "logs",
			ProtobufFieldName:  "logs",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		{
			Index:              5,
			Name:               "rolled_back_at",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "TIMESTAMPTZ",
			DatabaseTypePretty: "TIMESTAMPTZ",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMPTZ",
			ColumnLength:       -1,
			GoFieldName:        "RolledBackAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "rolled_back_at",
			ProtobufFieldName:  "rolled_back_at",
			ProtobufType:       "uint64",
			ProtobufPos:        6,
		},

		{
			Index:              6,
			Name:               "started_at",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TIMESTAMPTZ",
			DatabaseTypePretty: "TIMESTAMPTZ",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMPTZ",
			ColumnLength:       -1,
			GoFieldName:        "StartedAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "started_at",
			ProtobufFieldName:  "started_at",
			ProtobufType:       "uint64",
			ProtobufPos:        7,
		},

		{
			Index:              7,
			Name:               "applied_steps_count",
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
			GoFieldName:        "AppliedStepsCount",
			GoFieldType:        "int32",
			JSONFieldName:      "applied_steps_count",
			ProtobufFieldName:  "applied_steps_count",
			ProtobufType:       "int32",
			ProtobufPos:        8,
		},
	},
}

// TableName sets the insert table name for this struct type
func (p *PrismaMigrations) TableName() string {
	return "_prisma_migrations"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (p *PrismaMigrations) BeforeSave(*gorm.DB) error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (p *PrismaMigrations) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (p *PrismaMigrations) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (p *PrismaMigrations) TableInfo() *TableInfo {
	return _prisma_migrationsTableInfo
}
