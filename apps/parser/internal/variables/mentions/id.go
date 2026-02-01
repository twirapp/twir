package mentions

import (
	"context"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var ID = &types.Variable{
	Name: "mentions.id",
	Description: lo.ToPtr(
		`ID of mentioned user. Use $(mentions.id|N) to get Nth mention, default is 0. Use "all" to get comma separated mentions ids.`,
	),
	CommandsOnly:        false,
	CanBeUsedInRegistry: true,
	Visible:             lo.ToPtr(true),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		index := 0

		if variableData.Params != nil && *variableData.Params != "" {
			if *variableData.Params == "all" {
				var ids []string
				for _, mention := range parseCtx.Mentions {
					ids = append(ids, mention.UserId)
				}
				return &types.VariableHandlerResult{
					Result: strings.Join(ids, ","),
				}, nil
			} else {
				parsedIndex, err := strconv.Atoi(*variableData.Params)
				if err == nil && parsedIndex >= 0 {
					index = parsedIndex
				}
			}
		}

		if len(parseCtx.Mentions) == 0 || index >= len(parseCtx.Mentions) {
			return &types.VariableHandlerResult{
				Result: "",
			}, nil
		}

		return &types.VariableHandlerResult{
			Result: parseCtx.Mentions[index].UserId,
		}, nil
	},
}
