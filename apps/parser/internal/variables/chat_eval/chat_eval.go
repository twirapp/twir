package chat_eval

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
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

		res, err := parseCtx.Services.Executron.ExecuteUserCode(ctx, "javascript", script)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = "Probably you're doing some suspicious things or wrote wrong code."
			return result, nil
		}

		result.Result = res
		return result, nil
	},
}
