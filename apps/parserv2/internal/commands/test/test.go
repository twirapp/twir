package testcommand

import "tsuwari/parser/internal/types"

var desc = "test"
var module = "CHANNEL"

var Command = types.DefaultCommand{
	Command: types.Command{
		Name:        "",
		Description: &desc,
		Permission:  "MODERATOR",
		// Module: &module,
		Visible: true,
	},
	Handler: Handler,
}

func Handler(data types.VariableHandlerParams) {
	return
}
