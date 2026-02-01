package mentions

import (
	"context"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var Login = &types.Variable{
	Name: "mentions.login",
	Description: lo.ToPtr(
		`Login (username) of mentioned user. Use $(mentions.login|N) to get Nth mention, default is 0. Use "all" to get comma separated mentions logins.`,
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
				var logins []string
				for _, mention := range parseCtx.Mentions {
					logins = append(logins, mention.UserLogin)
				}
				return &types.VariableHandlerResult{
					Result: strings.Join(logins, ","),
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
			Result: parseCtx.Mentions[index].UserLogin,
		}, nil
	},
}
