package model

type ChannelRoleUser struct {
	ID     string `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID string `gorm:"column:userId;type:uuid;" json:"userId"`
	RoleID string `gorm:"column:roleId;type:uuid;" json:"-"`

	Role *ChannelRole `gorm:"foreignKey:RoleID" json:"-"`
	User *Users       `gorm:"foreignKey:UserID" json:"-"`
}

func (ChannelRoleUser) TableName() string {
	return "channels_roles_users"
}
