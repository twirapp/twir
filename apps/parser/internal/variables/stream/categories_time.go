package stream

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/pkg/helpers"
	"github.com/twirapp/twir/libs/repositories/channels_info_history/model"

	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
)

var CategoryTime = &types.Variable{
	Name:                "stream.category.time",
	Description:         lo.ToPtr("Time of current category, example: 1h 30m 22s"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		if parseCtx.ChannelStream != nil {
			result.Result = "Offline or error on getting category"
			return &result, nil
		}

		history, err := parseCtx.Services.ChannelsInfoHistoryRepo.GetMany(
			ctx,
			channelsinfohistory.GetManyInput{
				ChannelID: parseCtx.Channel.ID,
				After:     parseCtx.ChannelStream.StartedAt,
				Limit:     100,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Err:     err,
				Message: "Cannot get history of categories",
			}
		}

		if len(history) == 0 {
			result.Result = "No history recorded"
			return &result, nil
		}

		var category *model.ChannelInfoHistory
		for _, item := range history {
			if item.Category == parseCtx.ChannelStream.GameName {
				category = &item
				break
			}
		}

		if category == nil {
			result.Result = "No history recorded"
			return &result, nil
		}

		result.Result = helpers.Duration(
			category.CreatedAt,
			&helpers.DurationOpts{
				UseUtc: true,
				Hide:   helpers.DurationOptsHide{Years: true, Months: true, Days: true},
			},
		)

		return &result, nil
	},
}

// var CategoriesTime = &types.Variable{
// 	Name:                "stream.categories.time",
// 	Description:         lo.ToPtr("Time of played categories on current stream, example: Dota 2 × 1h 30m 22s · League of Legends × 1h 20m 10s"),
// 	CanBeUsedInRegistry: true,
// 	Handler: func(
// 		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
// 	) (*types.VariableHandlerResult, error) {
// 		result := types.VariableHandlerResult{}
//
// 		stream := parseCtx.Cacher.GetChannelStream(ctx)
// 		if stream == nil {
// 			result.Result = "Offline or error on getting category"
// 			return &result, nil
// 		}
//
// 		history, err := parseCtx.Services.ChannelsInfoHistoryRepo.GetMany(
// 			ctx,
// 			channelsinfohistory.GetManyInput{
// 				ChannelID: parseCtx.Channel.ID,
// 				After:     stream.StartedAt,
// 				Limit:     100,
// 			},
// 		)
// 		if err != nil {
// 			return nil, &types.CommandHandlerError{
// 				Err:     err,
// 				Message: "Cannot get history of categories",
// 			}
// 		}
//
// 		if len(history) == 0 {
// 			result.Result = "No history recorded"
// 			return &result, nil
// 		}
//
// 		categoriesStrings := make([]string, 0, len(history))
// 		for _, item := range history {
// 			categoriesStrings = append(
// 				categoriesStrings,
// 				fmt.Sprintf("%s × %s", item.Category, "jere"),
// 			)
// 		}
//
// 		slices.Reverse(categoriesStrings)
//
// 		result.Result = strings.Join(categoriesStrings, " => ")
//
// 		return &result, nil
// 	},
// }
