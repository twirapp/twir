package chat_eval

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/libs/grpc/eval"
)

var ChatEval = &types.Variable{
	Name:         "chatEval",
	Description:  lo.ToPtr("Evaluate custom script from chat"),
	Example:      lo.ToPtr("chatEval"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		if parseCtx.Text == nil || *parseCtx.Text == "" {
			return result, nil
		}

		script := fmt.Sprintf(`return %s`, *parseCtx.Text)

		req, err := parseCtx.Services.GrpcClients.Eval.Process(
			ctx,
			&eval.Evaluate{
				Script: script,
			},
		)

		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = "Probably you're doing some suspicious things."
			return result, nil
		}

		result.Result = lo.Substring(req.Result, 0, 500)
		return result, nil
	},
}
