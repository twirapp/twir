package predictions

import (
	"context"
	"strings"

	"github.com/guregu/null"
	"github.com/kvizyx/twitchy/helix"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/twitch"
)

const (
	startPredictionDuration    = "duration"
	startPredictionArgVariants = "variants"
	startPredictionArgTitle    = "title"
)

var Start = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "prediction start",
		Description: null.StringFrom("Start prediction. Example usage: !prediction start 100 | Will we win this game? | Yes, win / No, lose"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "PREDICTIONS",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	ArgsDelimiter:     " | ",
	Args: []command_arguments.Arg{
		command_arguments.Int{
			Name: startPredictionDuration,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Predictions.Hints.StartPredictionDuration,
				)
			},
		},
		command_arguments.String{
			Name: startPredictionArgTitle,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Predictions.Hints.StartPredictionArgTitle,
				)
			},
		},
		command_arguments.String{
			Name: startPredictionArgVariants,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Predictions.Hints.StartPredictionArgVariants,
				)
			},
		},
	},
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

		variantsArg := parseCtx.ArgsParser.Get(startPredictionArgVariants).String()
		durationArg := parseCtx.ArgsParser.Get(startPredictionDuration).Int()
		titleArg := parseCtx.ArgsParser.Get(startPredictionArgTitle).String()

		parsedVariants := strings.Split(variantsArg, " / ")
		outcomes := make([]helix.PredictionOutcome, 0, len(parsedVariants))
		for _, variant := range parsedVariants {
			outcomes = append(
				outcomes, helix.PredictionOutcome{
					Title: variant,
				},
			)
		}

		_, err = twitchClient.Predictions.CreatePrediction(
			ctx,
			helix.CreatePredictionRequest{
				BroadcasterID:    parseCtx.Channel.ID,
				Title:            titleArg,
				Outcomes:         outcomes,
				PredictionWindow: durationArg,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Predictions.Errors.CannotCreate),
				Err:     err,
			}
		}
		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(ctx, locales.Translations.Commands.Predictions.Info.Started),
			},
		}, nil
	},
}
