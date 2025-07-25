package emotes

import (
	"context"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

type sevenTVEmote struct {
	Name string `json:"name"`
}

type sevenUserTvResponse struct {
	EmoteSet *struct {
		Emotes []sevenTVEmote `json:"emotes"`
	} `json:"emote_set"`
}

var SevenTv = &types.Variable{
	Name:                "emotes.7tv",
	Description:         lo.ToPtr("Emotes of channel from https://7tv.app"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		var response sevenUserTvResponse

		_, err := req.R().
			SetContext(ctx).
			SetSuccessResult(&response).
			Get("https://7tv.io/v3/users/twitch/" + parseCtx.Channel.ID)

		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return result, nil
		}

		if response.EmoteSet == nil {
			return result, nil
		}

		mapped := lo.Map(
			response.EmoteSet.Emotes, func(e sevenTVEmote, _ int) string {
				return e.Name
			},
		)

		result.Result = strings.Join(mapped, " ")

		return result, nil
	},
}
