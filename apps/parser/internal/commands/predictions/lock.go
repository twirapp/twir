package predictions

import (
	"context"
	"errors"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
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

		cancelResp, err := twitchClient.EndPrediction(
			&helix.EndPredictionParams{
				BroadcasterID: parseCtx.Channel.ID,
				ID:            currentRunedPrediction.ID,
				Status:        "LOCKED",
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
				"✅ Prediction locked",
			},
		}, nil
	},
}
