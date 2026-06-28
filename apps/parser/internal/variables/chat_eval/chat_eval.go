package chat_eval

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	executronBus "github.com/twirapp/twir/libs/bus-core/executron"
	"github.com/twirapp/twir/libs/i18n"
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

		req, err := parseCtx.Services.Bus.Executron.Execute.Request(
			ctx,
			executronBus.ExecuteRequest{
				ChannelId: parseCtx.Channel.ID,
				Language:  "javascript",
				Code:      script,
			},
		)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.ChatEval.Info.WrongCode)
			return result, nil
		}

		var res string
		if req.Data.Result != "" {
			if utf8.RuneCountInString(req.Data.Result) > 474 {
				res = req.Data.Result[:474] + "..."
			} else {
				res = req.Data.Result
			}
		} else if req.Data.Error != "" {
			res = req.Data.Error
		}

		result.Result = res
		return result, nil
	},
}
