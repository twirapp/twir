package emotes

import (
	"context"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

type sevenTVEmote struct {
	Name string `json:"name"`
}

var SevenTv = &types.Variable{
	Name:                "emotes.7tv",
	Description:         lo.ToPtr("Emotes of channel from https://7tv.app"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		var data []*sevenTVEmote

		_, err := req.R().
			SetContext(ctx).
			SetSuccessResult(&data).
			Get("https://api.7tv.app/v2/users/" + parseCtx.Channel.ID + "/emotes")

		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return result, nil
		}

		mapped := lo.Map(
			data, func(e *sevenTVEmote, _ int) string {
				return e.Name
			},
		)

		result.Result = strings.Join(mapped, " ")

		return result, nil
	},
}
