package top

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

type emote struct {
	Emote string
	Count int

	UserID string
}

var Emotes = &types.Variable{
	Name:                "top.emotes",
	Description:         lo.ToPtr("Top used emotes"),
	Example:             lo.ToPtr("top.emotes|10"),
	CanBeUsedInRegistry: true,
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

		emotes, err := parseCtx.Services.ChannelEmotesUsagesRepo.GetEmotesStatistics(
			ctx,
			channelsemotesusagesrepository.GetEmotesStatisticsInput{
				ChannelID: parseCtx.Channel.ID,
				PerPage:   limit,
				Sort:      channelsemotesusagesrepository.SortDesc,
			},
		)
		if err != nil {
			return nil, err
		}

		mappedTop := lo.Map(
			emotes, func(e model.EmoteStatistic, _ int) string {
				return fmt.Sprintf(
					"%s × %v",
					e.EmoteName,
					e.TotalUsages,
				)
			},
		)

		result.Result = strings.Join(mappedTop, " ")
		return result, nil
	},
}
