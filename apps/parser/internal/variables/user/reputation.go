package user

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var Reputation = &types.Variable{
	Name:         "user.reputation",
	Description:  lo.ToPtr("User reputation"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		dbUser := parseCtx.Cacher.GetGbUserStats(ctx)
		if dbUser != nil {
			result.Result = strconv.FormatInt(dbUser.Reputation, 10)
		} else {
			result.Result = "0"
		}

		return &result, nil
	},
}
