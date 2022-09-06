package userage

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "user.age",
	Description: lo.ToPtr("User account age"),
	Handler: func(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		user := ctx.GetTwitchUser()
		if user == nil {
			result.Result = "error on getting user"
		} else {
			result.Result = helpers.Duration(user.CreatedAt.Time)
		}

		return &result, nil
	},
}
