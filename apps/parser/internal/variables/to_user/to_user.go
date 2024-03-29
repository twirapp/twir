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

		var mentions []string
		for _, m := range parseCtx.Mentions {
			mentions = append(mentions, "@"+m.UserName)
		}

		if mentions == nil {
			mentions = append(mentions, "@"+parseCtx.Sender.DisplayName)
		}

		result.Result = strings.Join(mentions, ", ")

		return result, nil
	},
}
