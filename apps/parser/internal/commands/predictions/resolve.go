package predictions

import (
	"context"
	"errors"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
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
			Hint: "variant number, for example: 1,2,3,4,5",
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
				Message: "cannot create twitch client",
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
				Message: "cannot get current prediction",
				Err:     err,
			}
		}
		if currentPredictionReq.ErrorMessage != "" {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf(
					"cannot get current prediction: %s",
					currentPredictionReq.ErrorMessage,
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
				Message: "no prediction runed",
			}
		}

		if len(currentRunedPrediction.Outcomes) < variant {
			return nil, &types.CommandHandlerError{
				Message: "no prediction variant",
			}
		}

		if variant > len(currentRunedPrediction.Outcomes) {
			return nil, &types.CommandHandlerError{
				Message: "no prediction variant",
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
				Message: "cannot cancel prediction",
				Err:     err,
			}
		}

		if cancelResp.ErrorMessage != "" {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf(
					"cannot cancel prediction: %s",
					cancelResp.ErrorMessage,
				),
				Err: errors.New(cancelResp.ErrorMessage),
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				"✅ Prediction resolved",
			},
		}, nil
	},
}
