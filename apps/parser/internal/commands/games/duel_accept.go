package games

import (
	"context"
	"errors"
	"math/rand"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"gorm.io/gorm"
)

type duelUserForTimeout struct {
	ID    string
	IsMod bool
	Name  string
}

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
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, nil
			}

			return nil, &types.CommandHandlerError{
				Message: "cannot get duel channel settings",
				Err:     err,
			}
		}

		currentDuel, err := handler.getUserCurrentDuel(ctx, parseCtx.Sender.ID)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get sender current duel",
				Err:     err,
			}
		}
		if currentDuel == nil {
			return &types.CommandsHandlerResult{
				Result: []string{"you are not participate in any duel"},
			}, nil
		}

		if currentDuel.TargetID.String != parseCtx.Sender.ID {
			return &types.CommandsHandlerResult{
				Result: []string{},
			}, nil
		}

		_, err = handler.getDbChannel(ctx)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get db channel",
				Err:     err,
			}
		}

		_, err = handler.createHelixClient()
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create broadcaster twitch client",
				Err:     err,
			}
		}

		randomedNumber := rand.Intn(100)

		var usersForTimeout []duelUserForTimeout
		var resultMessage string
		var loserId null.String

		if settings.BothDiePercent > 0 && randomedNumber <= int(settings.BothDiePercent) {
			usersForTimeout = append(
				usersForTimeout,
				duelUserForTimeout{
					ID:    currentDuel.SenderID.String,
					IsMod: currentDuel.SenderModerator,
					Name:  currentDuel.SenderLogin,
				},
				duelUserForTimeout{
					ID:    currentDuel.TargetID.String,
					IsMod: currentDuel.TargetModerator,
					Name:  currentDuel.TargetLogin,
				},
			)

			resultMessage = settings.BothDieMessage
			resultMessage = strings.ReplaceAll(resultMessage, "{initiator}", currentDuel.SenderLogin)
			resultMessage = strings.ReplaceAll(resultMessage, "{target}", currentDuel.TargetLogin)
		} else {
			remainderNumber := 100 - int(settings.BothDiePercent)
			var loser string
			var winner string
			if randomedNumber <= remainderNumber/2 {
				usersForTimeout = append(
					usersForTimeout,
					duelUserForTimeout{
						ID:    currentDuel.SenderID.String,
						IsMod: currentDuel.SenderModerator,
						Name:  currentDuel.SenderLogin,
					},
				)
				loser = currentDuel.SenderLogin
				winner = currentDuel.TargetLogin
				loserId = currentDuel.SenderID
			} else {
				usersForTimeout = append(
					usersForTimeout,
					duelUserForTimeout{
						ID:    currentDuel.TargetID.String,
						IsMod: currentDuel.TargetModerator,
						Name:  currentDuel.TargetLogin,
					},
				)
				loser = currentDuel.TargetLogin
				winner = currentDuel.SenderLogin
				loserId = currentDuel.TargetID
			}

			resultMessage = settings.ResultMessage
			resultMessage = strings.ReplaceAll(resultMessage, "{loser}", loser)
			resultMessage = strings.ReplaceAll(resultMessage, "{winner}", winner)
		}

		for _, user := range usersForTimeout {
			err = handler.timeoutUser(ctx, settings, user.ID, user.IsMod)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "cannot timeout user",
					Err:     err,
				}
			}
		}

		err = handler.saveResult(ctx, *currentDuel, settings, loserId)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot save duel result",
				Err:     err,
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{resultMessage},
		}, nil
	},
}
