package userwatched

import (
	"fmt"
	"time"

	types "github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:         "user.watched",
	Description:  lo.ToPtr("User watched time"),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		dbUser := ctx.GetGbUser()

		var watched int64 = 0

		if dbUser != nil {
			watched = dbUser.Watched
		}

		watchedD := time.Duration(watched) * time.Millisecond

		result.Result = fmt.Sprintf("%.1fh", watchedD.Hours())

		return &result, nil
	},
}
