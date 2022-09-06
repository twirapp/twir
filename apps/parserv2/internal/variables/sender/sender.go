package sender

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "sender",
	Description: lo.ToPtr("Username of user, who sended message"),
	Handler: func(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{Result: ctx.Context.SenderName}

		return &result, nil
	},
}
