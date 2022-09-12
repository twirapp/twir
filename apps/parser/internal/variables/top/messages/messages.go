package messages

import (
	"fmt"
	"strconv"
	"strings"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables/top"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "top.messages",
	Description: lo.ToPtr("Top users by messages"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}
		var page int = 1

		if ctx.Text != nil {
			p, err := strconv.Atoi(*ctx.Text)
			if err != nil {
				page = p
			}

			if page <= 0 {
				page = 1
			}
		}

		topUsers := top.GetTop(ctx, ctx.ChannelId, "messages", &page)

		if topUsers == nil || len(*topUsers) == 0 {
			return result, nil
		}

		mappedTop := lo.Map(*topUsers, func(user *top.UserStats, idx int) string {
			return fmt.Sprintf(
				"%v. %s — %v",
				(idx+1)+(page-1)*10,
				user.UserName,
				user.Value,
			)
		})

		result.Result = strings.Join(mappedTop, ", ")
		return result, nil
	},
}
