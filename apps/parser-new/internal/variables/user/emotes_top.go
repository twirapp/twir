package user

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"
)

type userEmotesTopEmote struct {
	Emote string
	Count int
}

var EmotesTop = &types.Variable{
	Name:         "user.top.emotes",
	Description:  lo.ToPtr("User top used emotes"),
	Example:      lo.ToPtr("user.top.emotes|10"),
	CommandsOnly: true,
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
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

		emotes := []userEmotesTopEmote{}
		err := parseCtx.Services.Gorm.
			Raw(`SELECT emote, COUNT(*)
				FROM channels_emotes_usages
				WHERE "channelId" = ? AND "userId" = ?
				Group By emote
				Order By COUNT(*)
				DESC LIMIT ?
			`, parseCtx.Channel.ID, parseCtx.Sender.ID, limit).
			Scan(&emotes).
			Error

		if err != nil {
			return nil, err
		}

		mappedTop := lo.Map(emotes, func(e userEmotesTopEmote, _ int) string {
			return fmt.Sprintf(
				"%s × %v",
				e.Emote,
				e.Count,
			)
		})

		result.Result = strings.Join(mappedTop, " ")
		return result, nil
	},
}
