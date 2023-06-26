package user

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var Watched = &types.Variable{
	Name:         "user.watched",
	Description:  lo.ToPtr("User watched time"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		dbUser := parseCtx.Cacher.GetGbUserStats(ctx)

		var watched int64 = 0

		if dbUser != nil {
			watched = dbUser.Watched
		}

		watchedD := time.Duration(watched) * time.Millisecond

		result.Result = fmt.Sprintf("%.1fh", watchedD.Hours())

		return &result, nil
	},
}
