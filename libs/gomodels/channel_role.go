package model

type ChannelRoleEnum string

const (
	ChannelRoleTypeBroadcaster ChannelRoleEnum = "BROADCASTER"
	ChannelRoleTypeModerator   ChannelRoleEnum = "MODERATOR"
	ChannelRoleTypeSubscriber  ChannelRoleEnum = "SUBSCRIBER"
	ChannelRoleTypeVip         ChannelRoleEnum = "VIP"
	//ChannelRoleTypeFollower    ChannelRoleEnum = "FOLLOWER"
	//ChannelRoleTypeViewer ChannelRoleEnum = "VIEWER"
	ChannelRoleTypeCustom ChannelRoleEnum = "CUSTOM"
)

type ChannelRole struct {
	ID        string          `gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	ChannelID string          `gorm:"column:channelId;type:uuid;"`
	Name      string          `gorm:"column:name;type:text;"`
	Type      ChannelRoleEnum `gorm:"column:type;type:text;"`
	System    bool            `gorm:"column:system;type:boolean;default:false"`

	Channel     *Channels                `gorm:"foreignKey:ChannelID"`
	Permissions []*ChannelRolePermission `gorm:"foreignKey:RoleID"`
	Users       []*ChannelRoleUser       `gorm:"foreignKey:RoleID"`
}

func (ChannelRole) TableName() string {
	return "channel_roles"
}
