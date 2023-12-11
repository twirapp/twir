package games

import (
	"context"
	"math/rand"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
)

var DuelAccept = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "duel accept",
		Description: null.StringFrom("Accept a duel with another user!"),
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
		handler := &duelHandler{parseCtx: parseCtx}

		settings, err := handler.getChannelSettings(ctx)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get duel channel settings",
				Err:     err,
			}
		}

		cachedData, err := handler.getSenderCurrentDuel(ctx)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get sender current duel",
				Err:     err,
			}
		}
		if cachedData.SenderID == "" {
			return &types.CommandsHandlerResult{
				Result: []string{"you are not participate in any duel"},
			}, nil
		}

		dbChannel, err := handler.getDbChannel(ctx)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get db channel",
				Err:     err,
			}
		}

		twitchClient, err := handler.createHelixClient(ctx)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create broadcaster twitch client",
				Err:     err,
			}
		}

		randomedNumber := rand.Intn(100)
		if settings.BothDiePercent > 0 && randomedNumber <= int(settings.BothDiePercent) {
			err = handler.timeoutUser(
				cachedData, dbChannel, settings, cachedData.SenderID, cachedData.IsSenderModerator,
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "cannot timeout user",
					Err:     err,
				}
			}

			err = handler.timeoutUser(
				cachedData, dbChannel, settings, cachedData.TargetID, cachedData.IsTargetModerator,
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "cannot timeout user",
					Err:     err,
				}
			}

			return &types.CommandsHandlerResult{
				Result: []string{settings.BothDieMessage},
			}, nil
		}

		remainderNumber := 100 - int(settings.BothDiePercent)
		var userId string
		var isMod bool

		if randomedNumber <= remainderNumber/2 {
			userId = cachedData.SenderID
			isMod = cachedData.IsSenderModerator
		} else {
			userId = cachedData.TargetID
			isMod = cachedData.IsTargetModerator
		}

		err = handler.timeoutUser(cachedData, dbChannel, settings, userId, isMod)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot timeout user",
				Err:     err,
			}
		}

		if isMod {
			go func() {
				time.Sleep(time.Duration(settings.TimeoutSeconds+2) * time.Second)
				_, err = twitchClient.AddChannelModerator(
					&helix.AddChannelModeratorParams{
						BroadcasterID: parseCtx.Channel.ID,
						UserID:        userId,
					},
				)
			}()
		}

		return nil, nil
	},
}
