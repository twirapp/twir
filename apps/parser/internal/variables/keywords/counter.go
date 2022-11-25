package keywords

import (
	"fmt"
	"strconv"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
)

var Counter = types.Variable{
	Name:         "keywords.counter",
	Description:  lo.ToPtr("Show how many times keyword was used"),
	CommandsOnly: lo.ToPtr(true),
	Visible:      lo.ToPtr(false),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		if data.Params == nil {
			result.Result = "id is not provided"
			return result, nil
		}

		keyword := model.ChannelsKeywords{}
		err := ctx.Services.Db.
			Where(`"channelId" = ? AND "id" = ?`, ctx.ChannelId, data.Params).
			Find(&keyword).Error
		if err != nil {
			fmt.Println(err)
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
