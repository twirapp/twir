package model

import "github.com/lib/pq"

type ChannelRole struct {
	ID          string          `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	ChannelID   string          `gorm:"column:channelId;type:uuid;" json:"-"`
	Name        string          `gorm:"column:name;type:text;" json:"name"`
	Type        ChannelRoleEnum `gorm:"column:type;type:text;" json:"type"`
	System      bool            `gorm:"column:system;type:boolean;default:false" json:"system"`
	Permissions pq.StringArray  `gorm:"column:permissions;type:text[]" json:"permissions"`

	Channel *Channels          `gorm:"foreignKey:ChannelID" json:"-"`
	Users   []*ChannelRoleUser `gorm:"foreignKey:RoleID" json:"-"`
}

func (ChannelRole) TableName() string {
	return "channels_roles"
}

type ChannelRoleEnum string

func (c ChannelRoleEnum) String() string {
	return string(c)
}

const (
	ChannelRoleTypeBroadcaster ChannelRoleEnum = "BROADCASTER"
	ChannelRoleTypeModerator   ChannelRoleEnum = "MODERATOR"
	ChannelRoleTypeSubscriber  ChannelRoleEnum = "SUBSCRIBER"
	ChannelRoleTypeVip         ChannelRoleEnum = "VIP"
	ChannelRoleTypeCustom      ChannelRoleEnum = "CUSTOM"
)

type RolePermissionEnum string

func (r RolePermissionEnum) String() string {
	return string(r)
}

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
