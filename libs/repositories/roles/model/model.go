package model

import (
	"github.com/google/uuid"
)

type Role struct {
	ID                        uuid.UUID
	ChannelID                 string
	Name                      string
	Type                      ChannelRoleEnum
	Permissions               []string
	RequiredWatchTime         int64
	RequiredMessages          int32
	RequiredUsedChannelPoints int64
}

var RoleNil = Role{}

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
