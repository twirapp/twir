package counttime

import (
	"context"
	"time"

	"github.com/araddon/dateparse"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/apps/parser/pkg/helpers"
	"github.com/twirapp/twir/libs/i18n"
)

var CountDown = &types.Variable{
	Name:         "countdown",
	Description:  lo.ToPtr("Shows countdown to date, support time."),
	CommandsOnly: true,
	Example:      lo.ToPtr("countdown|May 8, 2050 5:57:51 PM +0300"),
	Links: []types.VariableLink{
		{
			Name: "Supported formats",
			Href: "https://github.com/araddon/dateparse#extended-example",
		},
	},
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		if variableData.Params == nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Countdown.Errors.NotPassedParams)
			return result, nil
		}

		parsedTime, err := dateparse.ParseAny(*variableData.Params)
		if err != nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Countdown.Errors.ParseDate)
			return result, nil
		}

		result.Result = helpers.Duration(
			parsedTime,
			&helpers.DurationOpts{
				FromTime: time.Now(),
				Hide:     helpers.DurationOptsHide{},
			},
		)

		return result, nil
	},
}
