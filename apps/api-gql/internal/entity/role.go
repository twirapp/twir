package entity

import (
	"github.com/google/uuid"
)

type ChannelRole struct {
	ID                        uuid.UUID
	ChannelID                 string
	Name                      string
	Type                      ChannelRoleEnum
	Permissions               []string
	RequiredWatchTime         int64
	RequiredMessages          int32
	RequiredUsedChannelPoints int64
}

var ChannelRoleNil = ChannelRole{}

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

var AllChannelRoleTypeEnum = []ChannelRoleEnum{
	ChannelRoleTypeBroadcaster,
	ChannelRoleTypeModerator,
	ChannelRoleTypeSubscriber,
	ChannelRoleTypeVip,
	ChannelRoleTypeCustom,
}

type ChannelRoleUser struct {
	ID     uuid.UUID
	UserID string
	RoleID uuid.UUID
}

var ChannelRoleUserNil = ChannelRoleUser{}
