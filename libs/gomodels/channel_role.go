package model

type ChannelRoleEnum string

func (c ChannelRoleEnum) String() string {
	return string(c)
}

const (
	ChannelRoleTypeBroadcaster ChannelRoleEnum = "BROADCASTER"
	ChannelRoleTypeModerator   ChannelRoleEnum = "MODERATOR"
	ChannelRoleTypeSubscriber  ChannelRoleEnum = "SUBSCRIBER"
	ChannelRoleTypeVip         ChannelRoleEnum = "VIP"
	// ChannelRoleTypeFollower    ChannelRoleEnum = "FOLLOWER"
	// ChannelRoleTypeViewer ChannelRoleEnum = "VIEWER"
	ChannelRoleTypeCustom ChannelRoleEnum = "CUSTOM"
)

type ChannelRole struct {
	ID        string          `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	ChannelID string          `gorm:"column:channelId;type:uuid;" json:"-"`
	Name      string          `gorm:"column:name;type:text;" json:"name"`
	Type      ChannelRoleEnum `gorm:"column:type;type:text;" json:"type"`
	System    bool            `gorm:"column:system;type:boolean;default:false" json:"system"`

	Channel     *Channels                `gorm:"foreignKey:ChannelID" json:"-"`
	Permissions []*ChannelRolePermission `gorm:"foreignKey:RoleID" json:"permissions"`
	Users       []*ChannelRoleUser       `gorm:"foreignKey:RoleID" json:"-"`
}

func (ChannelRole) TableName() string {
	return "channels_roles"
}
