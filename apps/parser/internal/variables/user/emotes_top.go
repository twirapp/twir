package user

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	"github.com/twirapp/twir/libs/repositories/channels_emotes_usages/model"
)

var EmotesTop = &types.Variable{
	Name:         "user.top.emotes",
	Description:  lo.ToPtr("User top used emotes"),
	Example:      lo.ToPtr("user.top.emotes|10"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		limit := 10
		if variableData.Params != nil {
			newLimit, err := strconv.Atoi(*variableData.Params)
			if err == nil {
				limit = newLimit
			}
		}

		if limit > 50 {
			limit = 10
		}

		emotes, err := parseCtx.Services.ChannelEmotesUsagesRepo.GetUserMostUsedEmotes(
			ctx, channelsemotesusagesrepository.UserMostUsedEmotesInput{
				ChannelID: parseCtx.Channel.ID,
				UserID:    parseCtx.Sender.ID,
				Limit:     limit,
			},
		)
		if err != nil {
			return nil, err
		}

		mappedTop := lo.Map(
			emotes,
			func(e model.UserMostUsedEmote, _ int) string {
				return fmt.Sprintf(
					"%s × %v",
					e.Emote,
					e.Count,
				)
			},
		)

		result.Result = strings.Join(mappedTop, " ")
		return result, nil
	},
}
