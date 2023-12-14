package games

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"go.uber.org/zap"
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
		if cachedData == nil {
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

		twitchClient, err := handler.createHelixClient()
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create broadcaster twitch client",
				Err:     err,
			}
		}

		randomedNumber := rand.Intn(100)
		if settings.BothDiePercent > 0 && randomedNumber <= int(settings.BothDiePercent) {
			err = handler.timeoutUser(
				*cachedData, dbChannel, settings, cachedData.SenderID, cachedData.IsSenderModerator,
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "cannot timeout user",
					Err:     err,
				}
			}

			err = handler.timeoutUser(
				*cachedData, dbChannel, settings, cachedData.TargetID, cachedData.IsTargetModerator,
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "cannot timeout user",
					Err:     err,
				}
			}

			resultMessage := settings.BothDieMessage
			resultMessage = strings.ReplaceAll(resultMessage, "{initiator}", cachedData.SenderUserLogin)
			resultMessage = strings.ReplaceAll(resultMessage, "{target}", cachedData.TargetUserLogin)

			return &types.CommandsHandlerResult{
				Result: []string{resultMessage},
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

		err = handler.saveResult(ctx, *cachedData, dbChannel, settings, userId)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot save duel result",
				Err:     err,
			}
		}

		err = handler.timeoutUser(*cachedData, dbChannel, settings, userId, isMod)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot timeout user",
				Err:     err,
			}
		}

		if isMod {
			go func() {
				time.Sleep(time.Duration(settings.TimeoutSeconds+2) * time.Second)

				err := retry.Do(
					func() error {
						_, err = twitchClient.AddChannelModerator(
							&helix.AddChannelModeratorParams{
								BroadcasterID: parseCtx.Channel.ID,
								UserID:        userId,
							},
						)

						return err
					},
					retry.Attempts(5),
					retry.Delay(100*time.Millisecond),
				)

				if err != nil {
					parseCtx.Services.Logger.Error(
						"cannot add moderator",
						zap.Error(err),
					)
				}
			}()
		}

		var loserName string
		var winnerName string
		if userId == cachedData.SenderID {
			loserName = cachedData.SenderUserLogin
			winnerName = cachedData.TargetUserLogin
		} else {
			loserName = cachedData.TargetUserLogin
			winnerName = cachedData.SenderUserLogin
		}

		resultMessage := settings.ResultMessage
		resultMessage = strings.ReplaceAll(resultMessage, "{loser}", loserName)
		resultMessage = strings.ReplaceAll(resultMessage, "{winner}", winnerName)

		return &types.CommandsHandlerResult{
			Result: []string{resultMessage},
		}, nil
	},
}
