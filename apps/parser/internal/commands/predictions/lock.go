package predictions

import (
	"context"

	"github.com/guregu/null"
	"github.com/kvizyx/twitchy/helix"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/twitch"
)

var Lock = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "prediction lock",
		Description: null.StringFrom("Lock current prediction"),
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
			parseCtx.Channel.TwitchUserID,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.CannotCreateTwitch),
				Err:     err,
			}
		}

		currentPredictionReq, err := twitchClient.Predictions.GetPredictions(
			ctx,
			helix.GetPredictionsRequest{
				BroadcasterID: parseCtx.Channel.ID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.CannotCreateTwitch),
				Err:     err,
			}
		}
		var currentRunedPrediction *helix.Prediction
		for _, prediction := range currentPredictionReq.Data {
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

		_, err = twitchClient.Predictions.EndPrediction(
			ctx,
			helix.EndPredictionRequest{
				BroadcasterID: parseCtx.Channel.ID,
				ID:            currentRunedPrediction.ID,
				Status:        helix.PredictionStatusLocked,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Predictions.Errors.CannotCancel),
				Err:     err,
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(ctx, locales.Translations.Commands.Predictions.Info.Locked),
			},
		}, nil
	},
}
