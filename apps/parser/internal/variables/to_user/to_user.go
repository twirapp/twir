package to_user

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var ToUser = &types.Variable{
	Name:         "touser",
	Description:  lo.ToPtr("Mention user"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{
			Result: parseCtx.Sender.Name,
		}

		if parseCtx.Text != nil {
			result.Result = strings.ReplaceAll(*parseCtx.Text, "@", "")
		}

		return result, nil
	},
}
