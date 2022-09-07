package usermessages

import (
	"strconv"
	types "tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "user.messages",
	Description: lo.ToPtr("User messages"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		dbUser := ctx.GetGbUser()
		if dbUser != nil {
			result.Result = strconv.Itoa(int(dbUser.Messages))
		} else {
			result.Result = "0"
		}

		return &result, nil
	},
}
