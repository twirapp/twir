package model

import "fmt"

// Action CRUD actions
type Action int32

var (
	// Create action when record is created
	Create = Action(0)

	// RetrieveOne action when a record is retrieved from db
	RetrieveOne = Action(1)

	// RetrieveMany action when record(s) are retrieved from db
	RetrieveMany = Action(2)

	// Update action when record is updated in db
	Update = Action(3)

	// Delete action when record is deleted in db
	Delete = Action(4)

	// FetchDDL action when fetching ddl info from db
	FetchDDL = Action(5)

	tables map[string]*TableInfo
)

func init() {
	tables = make(map[string]*TableInfo)

	// tables["_prisma_migrations"] = _prisma_migrationsTableInfo
	tables["bots"] = botsTableInfo
	tables["channels_commands_usages"] = channels_commands_usagesTableInfo
	tables["channels_customvars"] = channels_customvarsTableInfo
	tables["channels_dashboard_access"] = channels_dashboard_accessTableInfo
	tables["channels_dota_accounts"] = channels_dota_accountsTableInfo
	tables["channels_integrations"] = channels_integrationsTableInfo
	tables["channels_permits"] = channels_permitsTableInfo
	tables["dota_game_modes"] = dota_game_modesTableInfo
	tables["dota_heroes"] = dota_heroesTableInfo
	tables["dota_matches"] = dota_matchesTableInfo
	tables["dota_matches_cards"] = dota_matches_cardsTableInfo
	tables["dota_matches_results"] = dota_matches_resultsTableInfo
	tables["dota_medals"] = dota_medalsTableInfo
	tables["integrations"] = integrationsTableInfo
	tables["notifications"] = notificationsTableInfo
	tables["notifications_messages"] = notifications_messagesTableInfo
	tables["tokens"] = tokensTableInfo
	tables["users"] = usersTableInfo
	tables["users_files"] = users_filesTableInfo
	tables["users_stats"] = users_statsTableInfo
	tables["users_viewed_notifications"] = users_viewed_notificationsTableInfo
}

// String describe the action
func (i Action) String() string {
	switch i {
	case Create:
		return "Create"
	case RetrieveOne:
		return "RetrieveOne"
	case RetrieveMany:
		return "RetrieveMany"
	case Update:
		return "Update"
	case Delete:
		return "Delete"
	case FetchDDL:
		return "FetchDDL"
	default:
		return fmt.Sprintf("unknown action: %d", int(i))
	}
}

// Model interface methods for database structs generated
type Model interface {
	TableName() string
	BeforeSave() error
	Prepare()
	Validate(action Action) error
	TableInfo() *TableInfo
}

// TableInfo describes a table in the database
type TableInfo struct {
	Name    string        `json:"name"`
	Columns []*ColumnInfo `json:"columns"`
}

// ColumnInfo describes a column in the database table
type ColumnInfo struct {
	Index              int    `json:"index"`
	GoFieldName        string `json:"go_field_name"`
	GoFieldType        string `json:"go_field_type"`
	JSONFieldName      string `json:"json_field_name"`
	ProtobufFieldName  string `json:"protobuf_field_name"`
	ProtobufType       string `json:"protobuf_field_type"`
	ProtobufPos        int    `json:"protobuf_field_pos"`
	Comment            string `json:"comment"`
	Notes              string `json:"notes"`
	Name               string `json:"name"`
	Nullable           bool   `json:"is_nullable"`
	DatabaseTypeName   string `json:"database_type_name"`
	DatabaseTypePretty string `json:"database_type_pretty"`
	IsPrimaryKey       bool   `json:"is_primary_key"`
	IsAutoIncrement    bool   `json:"is_auto_increment"`
	IsArray            bool   `json:"is_array"`
	ColumnType         string `json:"column_type"`
	ColumnLength       int64  `json:"column_length"`
	DefaultValue       string `json:"default_value"`
}

// GetTableInfo retrieve TableInfo for a table
func GetTableInfo(name string) (*TableInfo, bool) {
	val, ok := tables[name]
	return val, ok
}
