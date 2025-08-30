package stream

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
)

var Followers = &types.Variable{
	Name:                "stream.followers",
	Description:         lo.ToPtr("Followers on current stream"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		if parseCtx.ChannelStream == nil {
			result.Result = "0"
			return &result, nil
		}

		t := model.ChannelEventListItemTypeFollow
		count, err := parseCtx.Services.ChannelEventListsRepo.CountBy(
			ctx,
			channelseventslist.CountByInput{
				ChannelID:    &parseCtx.Channel.ID,
				CreatedAtGTE: &parseCtx.ChannelStream.StartedAt,
				Type:         &t,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Err:     err,
				Message: "Cannot count followers",
			}
		}

		result.Result = strconv.Itoa(int(count))

		return &result, nil
	},
}
