package model

import (
	rolecooldownmodel "github.com/twirapp/twir/libs/repositories/command_role_cooldown/model"
	commandmodel "github.com/twirapp/twir/libs/repositories/commands/model"
	groupmodel "github.com/twirapp/twir/libs/repositories/commands_group/model"
	responsemodel "github.com/twirapp/twir/libs/repositories/commands_response/model"
)

type CommandWithGroupAndResponses struct {
	commandmodel.Command
	Group         *groupmodel.Group                       `db:"group" json:"group"`
	Responses     []responsemodel.Response                `db:"responses" json:"responses"`
	RoleCooldowns []rolecooldownmodel.CommandRoleCooldown `db:"role_cooldowns" json:"role_cooldowns"`
}

var Nil = CommandWithGroupAndResponses{}
