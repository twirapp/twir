package keywords

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
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
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Keywords.Errors.IdNotProvided)
			return result, nil
		}

		var keyword *model.ChannelsKeywords
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND "id" = ?`, parseCtx.Channel.ID, variableData.Params).
			Find(&keyword).Error

		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)

			result.Result = i18n.GetCtx(ctx, locales.Translations.Errors.Generic.Internal)
			return result, nil
		}

		if keyword.ID == "" {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Keywords.Errors.NotFound)
			return result, nil
		}

		count := strconv.Itoa(int(keyword.Usages))
		result.Result = count

		return result, nil
	},
}
