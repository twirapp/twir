package model

type RolePermissionEnum string

const (
	RolePermissionCanAcessDashboard RolePermissionEnum = "CAN_ACCESS_DASHBOARD"

	RolePermissionUpdateChannelTitle    RolePermissionEnum = "UPDATE_CHANNEL_TITLE"
	RolePermissionUpdateChannelCategory RolePermissionEnum = "UPDATE_CHANNEL_CATEGORY"

	RolePermissionViewCommands   RolePermissionEnum = "VIEW_COMMANDS"
	RolePermissionManageCommands RolePermissionEnum = "MANAGE_COMMANDS"

	RolePermissionViewKeywords   RolePermissionEnum = "VIEW_KEYWORDS"
	RolePermissionManageKeywords RolePermissionEnum = "MANAGE_KEYWORDS"

	RolePermissionViewTimers   RolePermissionEnum = "VIEW_TIMERS"
	RolePermissionManageTimers RolePermissionEnum = "MANAGE_TIMERS"

	RolePermissionViewIntegrations   RolePermissionEnum = "VIEW_INTEGRATIONS"
	RolePermissionManageIntegrations RolePermissionEnum = "MANAGE_INTEGRATIONS"

	RolePermissionViewSongRequests   RolePermissionEnum = "VIEW_SONG_REQUESTS"
	RolePermissionManageSongRequests RolePermissionEnum = "MANAGE_SONG_REQUESTS"

	RolePermissionViewModeration   RolePermissionEnum = "VIEW_MODERATION"
	RolePermissionManageModeration RolePermissionEnum = "MANAGE_MODERATION"

	RolePermissionViewVariables   RolePermissionEnum = "VIEW_VARIABLES"
	RolePermissionManageVariables RolePermissionEnum = "MANAGE_VARIABLES"

	RolePermissionViewGreetings   RolePermissionEnum = "VIEW_GREETINGS"
	RolePermissionManageGreetings RolePermissionEnum = "MANAGE_GREETINGS"
)

type RoleFlag struct {
	ID   string             `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	Flag RolePermissionEnum `gorm:"column:flag;type:text;" json:"flag"`

	ChannelRolePermissions []*ChannelRolePermission `gorm:"foreignKey:FlagID" json:"-"`
}

func (RoleFlag) TableName() string {
	return "roles_flags"
}
