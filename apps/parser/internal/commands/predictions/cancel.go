package predictions

import (
	"context"
	"errors"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/twitch"
)

var Cancel = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "prediction cancel",
		Description: null.StringFrom("Cancel current prediction."),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "PREDICTIONS",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		twitchClient, err := twitch.NewUserClientWithContext(
			ctx,
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.CannotCreateTwitch),
				Err:     err,
			}
		}

		currentPredictionReq, err := twitchClient.GetPredictions(
			&helix.PredictionsParams{
				BroadcasterID: parseCtx.Channel.ID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Predictions.Errors.CannotGetCurrent),
				Err:     err,
			}
		}
		if currentPredictionReq.ErrorMessage != "" {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Predictions.Errors.CannotGetCurrentVar.
						SetVars(locales.KeysCommandsPredictionsErrorsCannotGetCurrentVarVars{Reason: currentPredictionReq.ErrorMessage}),
				),

				Err: errors.New(currentPredictionReq.ErrorMessage),
			}
		}

		var currentRunedPrediction *helix.Prediction
		for _, prediction := range currentPredictionReq.Data.Predictions {
			if prediction.Status == "LOCKED" || prediction.Status == "ACTIVE" {
				currentRunedPrediction = &prediction
				break
			}
		}

		if currentRunedPrediction == nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Predictions.Info.NoRuned),
			}
		}

		cancelResp, err := twitchClient.EndPrediction(
			&helix.EndPredictionParams{
				BroadcasterID: parseCtx.Channel.ID,
				ID:            currentRunedPrediction.ID,
				Status:        "CANCELED",
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Predictions.Errors.CannotCancel),
				Err:     err,
			}
		}

		if cancelResp.ErrorMessage != "" {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Predictions.Errors.CannotCancelVar.
						SetVars(locales.KeysCommandsPredictionsErrorsCannotCancelVarVars{Reason: cancelResp.ErrorMessage}),
				),
				Err: errors.New(cancelResp.ErrorMessage),
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(ctx, locales.Translations.Commands.Predictions.Info.Cancel),
			},
		}, nil
	},
}
