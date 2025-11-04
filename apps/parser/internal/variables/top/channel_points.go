package top

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var ChannelPoints = &types.Variable{
	Name:                "top.usedChannelPoints",
	Description:         lo.ToPtr("Top users by spent channel points"),
	Example:             lo.ToPtr("top.usedChannelPoints|10"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}
		var page = 1

		if parseCtx.Text != nil {
			p, err := strconv.Atoi(*parseCtx.Text)
			if err == nil {
				page = p
			}

			if page <= 0 {
				page = 1
			}
		}

		limit := 10
		if variableData.Params != nil {
			newLimit, err := strconv.Atoi(*variableData.Params)
			if err == nil {
				limit = newLimit
			}
		}

		if limit > 50 {
			limit = 10
		}

		topUsers := getTop(ctx, parseCtx, "usedChannelPoints", &page, limit)

		if topUsers == nil || len(topUsers) == 0 {
			return result, nil
		}

		mappedTop := lo.Map(
			topUsers, func(user *userStats, idx int) string {
				return fmt.Sprintf(
					"%s × %v",
					user.UserName,
					user.Value,
				)
			},
		)

		result.Result = strings.Join(mappedTop, " · ")
		return result, nil
	},
}
