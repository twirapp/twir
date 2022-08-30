package testcommand

import "tsuwari/parser/internal/types"

var desc = "test"
var module = "CHANNEL"

var Command = types.DefaultCommand{
	Command: types.Command{
		Name:        "test",
		Description: &desc,
		Permission:  "MODERATOR",
		Visible:     true,
	},
	Handler: Handler,
}

func Handler(data types.VariableHandlerParams) []string {
	return []string{"$(random|1-5000) 1", "$(random|6000-10000) 2"}
}
