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

type frankerFaceZEmote struct {
	Name string `json:"name"`
}

type frankerFaceZSet struct {
	Emoticons []*frankerFaceZEmote
}

type frankerFaceZResponse struct {
	Sets map[string]*frankerFaceZSet `json:"sets"`
}

var FrankerFaceZ = &types.Variable{
	Name:        "emotes.ffz",
	Description: lo.ToPtr("Emotes of channel from https://frankerfacez.com"),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.frankerfacez.com/v1/room/id/"+parseCtx.Channel.ID, nil)
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

		var data frankerFaceZResponse
		if err := json.Unmarshal(body, &data); err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return result, nil
		}

		var emotes []string

		for _, set := range data.Sets {
			mapped := lo.Map(
				set.Emoticons, func(e *frankerFaceZEmote, _ int) string {
					return e.Name
				},
			)

			emotes = append(emotes, mapped...)
		}

		result.Result = strings.Join(emotes, " ")

		return result, nil
	},
}
