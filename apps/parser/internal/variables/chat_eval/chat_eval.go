package chat_eval

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var ChatEval = &types.Variable{
	Name:                     "chatEval",
	Description:              lo.ToPtr("Evaluate custom script from chat"),
	Example:                  lo.ToPtr("chatEval"),
	DisableInCustomVariables: true,
	NotCachable:              true,
	CommandsOnly:             true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		if parseCtx.Text == nil || *parseCtx.Text == "" {
			return result, nil
		}

		script := fmt.Sprintf(`return %s`, *parseCtx.Text)

		req, err := parseCtx.Services.Executron.ExecuteUserCode(
			ctx,
			parseCtx.Channel.ID,
			"javascript",
			script,
		)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = "Probably you're doing some suspicious things or wrote wrong code."
			return result, nil
		}

		var res string
		if req.Result != "" {
			if utf8.RuneCountInString(req.Result) > 474 {
				res = req.Result[:474] + "..."
			} else {
				res = req.Result
			}
		} else if req.Error != "" {
			res = req.Error
		}

		result.Result = res
		return result, nil
	},
}
