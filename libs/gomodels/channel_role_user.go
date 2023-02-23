package model

type ChannelRoleUser struct {
	ID     string `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	UserID string `gorm:"column:userId;type:uuid;"`
	RoleID string `gorm:"column:roleId;type:uuid;"`

	Role *ChannelRole `gorm:"foreignKey:RoleID"`
	User *Users       `gorm:"foreignKey:UserID"`
}

func (ChannelRoleUser) TableName() string {
	return "channels_roles_users"
}
