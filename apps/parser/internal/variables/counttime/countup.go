package counttime

import (
	"context"
	"time"

	"github.com/araddon/dateparse"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/pkg/helpers"
)

var CountUp = &types.Variable{
	Name:         "countup",
	Description:  lo.ToPtr("Shows time passed from time."),
	CommandsOnly: true,
	Example:      lo.ToPtr("countup|Oct 5, 1998 5:57:51 PM +0300"),
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
			result.Result = "Have not passed params to variable. "
			return result, nil
		}

		parsedTime, err := dateparse.ParseAny(*variableData.Params)
		if err != nil {
			result.Result = "Cannot parse date"
			return result, nil
		}

		result.Result = helpers.Duration(
			time.Now(),
			&helpers.DurationOpts{
				FromTime: parsedTime,
				Hide:     helpers.DurationOptsHide{},
			},
		)

		return result, nil
	},
}
