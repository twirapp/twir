package predictions

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
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
			Hint: "120",
		},
		command_arguments.String{
			Name: startPredictionArgTitle,
			Hint: "Will we win this game?",
		},
		command_arguments.String{
			Name: startPredictionArgVariants,
			Hint: "Yes, win / No, lose",
		},
	},
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
				Message: "cannot create twitch client",
				Err:     err,
			}
		}

		variantsArg := parseCtx.ArgsParser.Get(startPredictionArgVariants).String()
		durationArg := parseCtx.ArgsParser.Get(startPredictionDuration).Int()
		titleArg := parseCtx.ArgsParser.Get(startPredictionArgTitle).String()

		parsedVariants := strings.Split(variantsArg, " / ")
		outcomes := make([]helix.PredictionChoiceParam, 0, len(parsedVariants))
		for _, variant := range parsedVariants {
			outcomes = append(
				outcomes, helix.PredictionChoiceParam{
					Title: variant,
				},
			)
		}

		createResp, err := twitchClient.CreatePrediction(
			&helix.CreatePredictionParams{
				BroadcasterID:    parseCtx.Channel.ID,
				Title:            titleArg,
				Outcomes:         outcomes,
				PredictionWindow: durationArg,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create prediction",
				Err:     err,
			}
		}
		if createResp.ErrorMessage != "" {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf(
					"cannot create prediction: %s",
					createResp.ErrorMessage,
				),
				Err: errors.New(createResp.ErrorMessage),
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				"✅ Prediction started",
			},
		}, nil
	},
}
