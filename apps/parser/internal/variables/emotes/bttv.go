package emotes

import (
	"context"
	"strings"

	"github.com/imroc/req/v3"
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

		var data *betterTTVResponse

		_, err := req.R().
			SetContext(ctx).
			SetSuccessResult(&data).
			Get("https://api.betterttv.net/3/cached/users/twitch/" + parseCtx.Channel.ID)

		if err != nil {
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
