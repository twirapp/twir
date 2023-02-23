package model

type RolePermissionEnum string

const (
	RolePermissionAdministrator RolePermissionEnum = "ADMINISTRATOR"

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

type RolePermission struct {
	ID         string             `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	Permission RolePermissionEnum `gorm:"column:permission;type:text;"`

	ChannelRolePermissions []*ChannelRolePermission `gorm:"foreignKey:PermissionID"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
