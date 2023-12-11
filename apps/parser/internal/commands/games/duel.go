package games

import (
	"context"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
)

var Duel = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "duel",
		Description: null.StringFrom("Start a duel with another user!"),
		Module:      "GAMES",
		IsReply:     false,
		Visible:     true,
		Enabled:     false,
		RolesIDS:    pq.StringArray{},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		handler := &duelHandler{
			parseCtx: parseCtx,
		}

		errorResult := types.CommandsHandlerResult{
			Result: []string{"Something went wrong, please try again later"},
		}

		settings, err := handler.getChannelSettings(ctx)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get duel channel settings",
				Err:     err,
			}
		}

		if !settings.Enabled {
			return &types.CommandsHandlerResult{
				Result: []string{},
			}, nil
		}

		if parseCtx.Text == nil || *parseCtx.Text == "" {
			return &types.CommandsHandlerResult{}, nil
		}

		dbChannel, err := handler.getDbChannel(ctx)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get db channel",
				Err:     err,
			}
		}

		_, err = handler.createHelixClient(ctx)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create broadcaster twitch client",
				Err:     err,
			}
		}

		targetUser, err := handler.getTwitchTargetUser()
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot find target user on twitch",
				Err:     err,
			}
		}

		if err = handler.validateTarget(ctx, targetUser, dbChannel); err != nil {
			errorResult.Result = []string{err.Error()}
			return &errorResult, nil
		}

		moderators, err := handler.getChannelModerators()
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get channel moderators",
				Err:     err,
			}
		}

		if err = handler.saveDuelDataToCache(ctx, targetUser, moderators, settings); err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot save duel data to cache",
				Err:     err,
			}
		}

		var result []string
		if settings.StartMessage != "" {
			startMessage := settings.StartMessage
			startMessage = strings.ReplaceAll(startMessage, "{target}", targetUser.Login)
			startMessage = strings.ReplaceAll(startMessage, "{sender}", parseCtx.Sender.Name)
			startMessage = strings.ReplaceAll(
				startMessage,
				"{accept_seconds}",
				fmt.Sprintf("%v", settings.SecondsToAccept),
			)

			result = append(result, startMessage)
		}

		return &types.CommandsHandlerResult{
			Result: result,
		}, nil
	},
}

func generateDuelRedisKey(channelID, senderID, targetID string) string {
	return fmt.Sprintf("channels:%v:commands:duel:%v:%v", channelID, senderID, targetID)
}
