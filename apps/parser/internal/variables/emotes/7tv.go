package emotes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

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

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://7tv.io/v3/users/twitch/"+parseCtx.Channel.ID, nil)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return result, nil
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return result, nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return result, nil
		}

		var response sevenUserTvResponse
		if err := json.Unmarshal(body, &response); err != nil {
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
