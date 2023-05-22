package user

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"
)

var Messages = &types.Variable{
	Name:         "user.messages",
	Description:  lo.ToPtr("User messages"),
	CommandsOnly: true,
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		dbUser := parseCtx.Cacher.GetGbUserStats(ctx)
		if dbUser != nil {
			result.Result = strconv.Itoa(int(dbUser.Messages))
		} else {
			result.Result = "0"
		}

		return &result, nil
	},
}
