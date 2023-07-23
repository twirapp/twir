package handler

import (
	"context"

	"github.com/dnsge/twitch-eventsub-bindings"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/grpc/generated/events"
	"go.uber.org/zap"
)

func convertOutCome(outcomes []eventsub_bindings.PredictionOutcome) []*events.PredictionInfo_OutCome {
	out := make([]*events.PredictionInfo_OutCome, 0, len(outcomes))

	for _, outcome := range outcomes {
		topPredictors := make([]*events.PredictionInfo_OutCome_TopPredictor, 0, len(outcome.TopPredictors))

		for _, predictor := range outcome.TopPredictors {
			won := uint64(predictor.ChannelPointsWon)
			topPredictors = append(
				topPredictors, &events.PredictionInfo_OutCome_TopPredictor{
					UserName:        predictor.UserLogin,
					UserDisplayName: predictor.UserName,
					UserId:          predictor.UserID,
					PointsUsed:      uint64(predictor.ChannelPointsUsed),
					PointsWin: lo.
						If(predictor.ChannelPointsWon > 0, &won).
						Else(nil),
				},
			)
		}

		out = append(
			out, &events.PredictionInfo_OutCome{
				Id:            outcome.ID,
				Title:         outcome.Title,
				Color:         outcome.Color,
				Users:         uint64(outcome.Users),
				ChannelPoints: uint64(outcome.ChannelPoints),
				TopPredictors: topPredictors,
			},
		)
	}
	return out
}

func (c *Handler) handleChannelPredictionBegin(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPredictionBegin,
) {
	zap.S().Infow(
		"Prediction begin",
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
		"title", event.Title,
		"outcomes", event.Outcomes,
	)

	outComes := convertOutCome(event.Outcomes)

	_, err := c.services.Grpc.Events.PredictionBegin(
		context.Background(),
		&events.PredictionBeginMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			UserName:        event.BroadcasterUserLogin,
			UserDisplayName: event.BroadcasterUserName,
			Info: &events.PredictionInfo{
				Title:    event.Title,
				Outcomes: outComes,
			},
		},
	)

	if err != nil {
		zap.S().Error(err)
	}
}

func (c *Handler) handleChannelPredictionProgress(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPredictionProgress,
) {
	zap.S().Infow(
		"Prediction progress",
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
		"title", event.Title,
		"outcomes", event.Outcomes,
	)

	outComes := convertOutCome(event.Outcomes)

	_, err := c.services.Grpc.Events.PredictionProgress(
		context.Background(),
		&events.PredictionProgressMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			UserName:        event.BroadcasterUserLogin,
			UserDisplayName: event.BroadcasterUserName,
			Info: &events.PredictionInfo{
				Title:    event.Title,
				Outcomes: outComes,
			},
		},
	)

	if err != nil {
		zap.S().Error(err)
	}
}

func (c *Handler) handleChannelPredictionLock(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPredictionLock,
) {
	zap.S().Infow(
		"Prediction lock",
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
		"title", event.Title,
		"outcomes", event.Outcomes,
	)

	outComes := convertOutCome(event.Outcomes)

	_, err := c.services.Grpc.Events.PredictionLock(
		context.Background(),
		&events.PredictionLockMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			UserName:        event.BroadcasterUserLogin,
			UserDisplayName: event.BroadcasterUserName,
			Info: &events.PredictionInfo{
				Title:    event.Title,
				Outcomes: outComes,
			},
		},
	)

	if err != nil {
		zap.S().Error(err)
	}
}

func (c *Handler) handleChannelPredictionEnd(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPredictionEnd,
) {
	zap.S().Infow(
		"Prediction end",
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
		"title", event.Title,
		"outcomes", event.Outcomes,
		"status", event.Status,
	)

	if event.Status != "resolved" {
		return
	}

	outComes := convertOutCome(event.Outcomes)

	_, err := c.services.Grpc.Events.PredictionEnd(
		context.Background(),
		&events.PredictionEndMessage{
			BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
			UserName:        event.BroadcasterUserLogin,
			UserDisplayName: event.BroadcasterUserName,
			Info: &events.PredictionInfo{
				Title:    event.Title,
				Outcomes: outComes,
			},
			WinningOutcomeId: event.WinningOutcomeID,
		},
	)

	if err != nil {
		zap.S().Error(err)
	}
}
