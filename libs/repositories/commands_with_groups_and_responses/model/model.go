package model

import (
	commandmodel "github.com/twirapp/twir/libs/repositories/commands/model"
	groupmodel "github.com/twirapp/twir/libs/repositories/commands_group/model"
	responsemodel "github.com/twirapp/twir/libs/repositories/commands_responses/model"
)

type CommandWithGroupAndResponses struct {
	Command   commandmodel.Command
	Group     *groupmodel.Group
	Responses []responsemodel.Response
}

var Nil = CommandWithGroupAndResponses{}
