package testcommand

import (
	"tsuwari/parser/internal/types"

	"github.com/samber/lo"
)

var Command = types.DefaultCommand{
	Command: types.Command{
		Name:        "test",
		Description: lo.ToPtr("test"),
		Permission:  "MODERATOR",
		Visible:     true,
		Module:      lo.ToPtr("CHANNEL"),
	},
	Handler: Handler,
}

func Handler(data types.VariableHandlerParams) []string {
	return []string{"$(random|1-5000) 1", "$(random|6000-10000) 2"}
}
