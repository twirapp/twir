package model

type ChannelRolePermission struct {
	ID           string `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	RoleID       string `gorm:"column:roleId;type:uuid;"`
	PermissionID string `gorm:"column:permissionId;type:uuid;"`

	Role       *ChannelRole    `gorm:"foreignKey:RoleID"`
	Permission *RolePermission `gorm:"foreignKey:PermissionID"`
}

func (ChannelRolePermission) TableName() string {
	return "channels_roles_permissions"
}
