package stream

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/i18n"
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

		channelID := parseCtx.Channel.DBChannelID
		if channelID == "" {
			channelID = parseCtx.Channel.ID
		}

		eventPlatform := parseCtx.Platform
		if eventPlatform == "" {
			eventPlatform = platform.PlatformTwitch
		}

		t := model.ChannelEventListItemTypeFollow
		count, err := parseCtx.Services.ChannelEventListsRepo.CountBy(
			ctx,
			channelseventslist.CountByInput{
				ChannelID:    &channelID,
				Platform:     &eventPlatform,
				CreatedAtGTE: &parseCtx.ChannelStream.StartedAt,
				Type:         &t,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Err:     err,
				Message: i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.CountFollowers),
			}
		}

		result.Result = strconv.Itoa(int(count))

		return &result, nil
	},
}
