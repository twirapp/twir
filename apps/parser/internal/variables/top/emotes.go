package top

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

type emote struct {
	Emote string
	Count int

	UserID string
}

var Emotes = &types.Variable{
	Name:        "top.emotes",
	Description: lo.ToPtr("Top used emotes"),
	Example:     lo.ToPtr("top.emotes|10"),
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

		emotes := []emote{}
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Raw(
				`SELECT emote, COUNT(*)
				FROM channels_emotes_usages
				WHERE "channelId" = ?
				Group By emote
				Order By COUNT(*)
				DESC LIMIT ?
			`, parseCtx.Channel.ID, limit,
			).
			Scan(&emotes).
			Error

		if err != nil {
			return nil, err
		}

		mappedTop := lo.Map(
			emotes, func(e emote, _ int) string {
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
