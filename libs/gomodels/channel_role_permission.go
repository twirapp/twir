package model

type ChannelRolePermission struct {
	ID     string `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	RoleID string `gorm:"column:roleId;type:uuid;" json:"-"`
	FlagID string `gorm:"column:flagId;type:uuid;" json:"-"`

	Role *ChannelRole `gorm:"foreignKey:RoleID" json:"-"`
	Flag *RoleFlag    `gorm:"foreignKey:FlagID" json:"flag"`
}

func (ChannelRolePermission) TableName() string {
	return "channels_roles_permissions"
}
