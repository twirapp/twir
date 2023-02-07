package watched

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	"github.com/satont/tsuwari/apps/parser/internal/variables/top"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"strconv"
	"strings"
	"time"
)

var Variable = types.Variable{
	Name:        "top.watched",
	Description: lo.ToPtr("Top users by watch time"),
	Example:     lo.ToPtr("top.watched|10"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}
		var page = 1

		if ctx.Text != nil {
			p, err := strconv.Atoi(*ctx.Text)
			if err != nil {
				page = p
			}

			if page <= 0 {
				page = 1
			}
		}

		limit := 10
		if data.Params != nil {
			newLimit, err := strconv.Atoi(*data.Params)
			if err == nil {
				limit = newLimit
			}
		}

		if limit > 50 {
			limit = 10
		}

		topUsers := top.GetTop(ctx, ctx.ChannelId, "watched", &page, limit)

		if topUsers == nil || len(topUsers) == 0 {
			return result, nil
		}

		mappedTop := lo.Map(topUsers, func(user *top.UserStats, idx int) string {
			duration := time.Duration(user.Value) * time.Millisecond
			return fmt.Sprintf(
				"%s %.1fh",
				user.UserName,
				duration.Hours(),
			)
		})

		result.Result = strings.Join(mappedTop, " · ")
		return result, nil
	},
}
