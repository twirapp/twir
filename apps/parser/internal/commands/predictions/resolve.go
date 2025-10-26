package predictions

import (
	"context"
	"errors"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/twitch"
)

const (
	predictionResolveOutcomeNum = "outcome"
)

var Resolve = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "prediction resolve",
		Description: null.StringFrom("Cancel current prediction."),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "PREDICTIONS",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.Int{
			Name: predictionResolveOutcomeNum,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Predictions.Hints.PredictionResolveOutcomeNum,
				)
			},
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		variant := parseCtx.ArgsParser.Get(predictionResolveOutcomeNum).Int()

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
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Predictions.Errors.CannotGetCurrent,
				),
				Err: err,
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

		if len(currentRunedPrediction.Outcomes) < variant {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Predictions.Errors.NoVariant),
			}
		}

		if variant > len(currentRunedPrediction.Outcomes) {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Predictions.Errors.NoVariant),
			}
		}

		foundOutcome := currentRunedPrediction.Outcomes[variant-1]

		cancelResp, err := twitchClient.EndPrediction(
			&helix.EndPredictionParams{
				BroadcasterID:    parseCtx.Channel.ID,
				ID:               currentRunedPrediction.ID,
				Status:           "RESOLVED",
				WinningOutcomeID: foundOutcome.ID,
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
				i18n.GetCtx(ctx, locales.Translations.Commands.Predictions.Info.Resolved),
			},
		}, nil
	},
}
