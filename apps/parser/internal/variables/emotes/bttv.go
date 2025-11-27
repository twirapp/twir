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

type betterTTVEmote struct {
	Code string `json:"code"`
}

type betterTTVResponse struct {
	ChannelEmotes []*betterTTVEmote `json:"channelEmotes"`
	SharedEmotes  []*betterTTVEmote `json:"sharedEmotes"`
}

var BetterTTV = &types.Variable{
	Name:        "emotes.bttv",
	Description: lo.ToPtr("Emotes of channel from https://betterttv.com/"),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.betterttv.net/3/cached/users/twitch/"+parseCtx.Channel.ID, nil)
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

		var data betterTTVResponse
		if err := json.Unmarshal(body, &data); err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return result, nil
		}

		emotes := make([]string, 0, len(data.ChannelEmotes)+len(data.SharedEmotes))

		mappedChannelEmotes := lo.Map(
			data.ChannelEmotes, func(e *betterTTVEmote, _ int) string {
				return e.Code
			},
		)
		mappedSharedEmotes := lo.Map(
			data.SharedEmotes, func(e *betterTTVEmote, _ int) string {
				return e.Code
			},
		)

		emotes = append(emotes, mappedChannelEmotes...)
		emotes = append(emotes, mappedSharedEmotes...)

		result.Result = strings.Join(emotes, " ")

		return result, nil
	},
}
