package keywords

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
)

var Counter = &types.Variable{
	Name:        "keywords.counter",
	Description: lo.ToPtr("Show how many times keyword was used"),
	Example:     lo.ToPtr("keywords.counter|id"),
	Visible:     lo.ToPtr(false),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		if variableData.Params == nil {
			result.Result = "id is not provided"
			return result, nil
		}

		var keyword *model.ChannelsKeywords
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND "id" = ?`, parseCtx.Channel.ID, variableData.Params).
			Find(&keyword).Error

		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)

			result.Result = "internal error"
			return result, nil
		}

		if keyword.ID == "" {
			result.Result = "keyword not found"
			return result, nil
		}

		count := strconv.Itoa(int(keyword.Usages))
		result.Result = count

		return result, nil
	},
}
