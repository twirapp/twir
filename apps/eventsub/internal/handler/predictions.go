package handler

import (
	"context"
	"log/slog"

	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/grpc/events"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func convertOutCome(outcomes []eventsub_bindings.PredictionOutcome) []*events.PredictionInfo_OutCome {
	out := make([]*events.PredictionInfo_OutCome, 0, len(outcomes))

	for _, outcome := range outcomes {
		topPredictors := make(
			[]*events.PredictionInfo_OutCome_TopPredictor,
			0,
			len(outcome.TopPredictors),
		)

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
	c.logger.Info(
		"Prediction begin",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("title", event.Title),
	)

	outComes := convertOutCome(event.Outcomes)

	_, err := c.eventsGrpc.PredictionBegin(
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
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}

func (c *Handler) handleChannelPredictionProgress(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPredictionProgress,
) {
	c.logger.Info(
		"Prediction progress",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("title", event.Title),
	)

	outComes := convertOutCome(event.Outcomes)

	_, err := c.eventsGrpc.PredictionProgress(
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
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}

func (c *Handler) handleChannelPredictionLock(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPredictionLock,
) {
	c.logger.Info(
		"Prediction lock",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("title", event.Title),
	)

	outComes := convertOutCome(event.Outcomes)

	_, err := c.eventsGrpc.PredictionLock(
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
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}

func (c *Handler) handleChannelPredictionEnd(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPredictionEnd,
) {
	c.logger.Info(
		"Prediction end",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("title", event.Title),
		slog.String("status", event.Status),
	)

	if event.Status != "resolved" {
		return
	}

	outComes := convertOutCome(event.Outcomes)

	_, err := c.eventsGrpc.PredictionEnd(
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
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
