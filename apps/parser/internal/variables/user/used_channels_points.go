package user

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var UsedChannelPoints = &types.Variable{
	Name:         "user.usedChannelPoints",
	Description:  lo.ToPtr("How many channel points user spent on channel"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		dbUser := parseCtx.Cacher.GetGbUserStats(ctx)

		if dbUser != nil {
			result.Result = strconv.Itoa(int(dbUser.UsedChannelPoints))
		} else {
			result.Result = "0"
		}

		return &result, nil
	},
}
